package redditkv

import (
	"testing"
)

func TestSetAndGet(t *testing.T) {
	mock := NewMockRedditAPI()
	client := NewWithAPI(mock, "testsubreddit")

	// Set a key
	err := client.Set("mykey", "myvalue")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Get the key
	value, err := client.Get("mykey")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if value.Value != "myvalue" {
		t.Errorf("Expected value 'myvalue', got '%s'", value.Value)
	}

	if len(value.Children) != 0 {
		t.Errorf("Expected no children, got %d", len(value.Children))
	}
}

func TestSetOverwrites(t *testing.T) {
	mock := NewMockRedditAPI()
	client := NewWithAPI(mock, "testsubreddit")

	// Set a key
	err := client.Set("mykey", "value1")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Overwrite the key
	err = client.Set("mykey", "value2")
	if err != nil {
		t.Fatalf("Set (overwrite) failed: %v", err)
	}

	// Get should return new value
	value, err := client.Get("mykey")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if value.Value != "value2" {
		t.Errorf("Expected value 'value2', got '%s'", value.Value)
	}

	// Should only have 1 post (the overwrite deleted the old one)
	if mock.GetPostCount() != 1 {
		t.Errorf("Expected 1 post after overwrite, got %d", mock.GetPostCount())
	}
}

func TestAppend(t *testing.T) {
	mock := NewMockRedditAPI()
	client := NewWithAPI(mock, "testsubreddit")

	// Set initial value
	err := client.Set("mykey", "root")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Append a sibling
	err = client.Append("mykey", "sibling", nil)
	if err != nil {
		t.Fatalf("Append failed: %v", err)
	}

	// Get and verify structure
	value, err := client.Get("mykey")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if value.Value != "root" {
		t.Errorf("Expected root value 'root', got '%s'", value.Value)
	}

	// The sibling should be a child of root (how we model top-level siblings)
	if len(value.Children) != 1 {
		t.Errorf("Expected 1 child, got %d", len(value.Children))
	}

	if value.Children[0].Value != "sibling" {
		t.Errorf("Expected child value 'sibling', got '%s'", value.Children[0].Value)
	}
}

func TestAppendWithPath(t *testing.T) {
	mock := NewMockRedditAPI()
	client := NewWithAPI(mock, "testsubreddit")

	// Set initial value
	err := client.Set("mykey", "root")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Append a child to root (path [0] means first comment)
	err = client.Append("mykey", "child", []int{0})
	if err != nil {
		t.Fatalf("Append with path failed: %v", err)
	}

	// Get and verify structure
	value, err := client.Get("mykey")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if value.Value != "root" {
		t.Errorf("Expected root value 'root', got '%s'", value.Value)
	}

	if len(value.Children) != 1 {
		t.Errorf("Expected 1 child, got %d", len(value.Children))
	}

	if value.Children[0].Value != "child" {
		t.Errorf("Expected child value 'child', got '%s'", value.Children[0].Value)
	}
}

func TestDelete(t *testing.T) {
	mock := NewMockRedditAPI()
	client := NewWithAPI(mock, "testsubreddit")

	// Set a key
	err := client.Set("mykey", "myvalue")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Delete the key
	err = client.Delete("mykey")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Get should fail
	_, err = client.Get("mykey")
	if err == nil {
		t.Error("Expected error getting deleted key")
	}

	if _, ok := err.(*KeyNotFoundError); !ok {
		t.Errorf("Expected KeyNotFoundError, got %T", err)
	}
}

func TestKeys(t *testing.T) {
	mock := NewMockRedditAPI()
	client := NewWithAPI(mock, "testsubreddit")

	// Set some keys
	_ = client.Set("key1", "value1")
	_ = client.Set("key2", "value2")
	_ = client.Set("key3", "value3")

	keys, err := client.Keys()
	if err != nil {
		t.Fatalf("Keys failed: %v", err)
	}

	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

	// Keys should contain all our keys (order not guaranteed)
	keySet := make(map[string]bool)
	for _, k := range keys {
		keySet[k] = true
	}

	for _, expected := range []string{"key1", "key2", "key3"} {
		if !keySet[expected] {
			t.Errorf("Expected key '%s' not found", expected)
		}
	}
}

func TestExists(t *testing.T) {
	mock := NewMockRedditAPI()
	client := NewWithAPI(mock, "testsubreddit")

	// Key shouldn't exist initially
	exists, err := client.Exists("mykey")
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if exists {
		t.Error("Expected key to not exist")
	}

	// Set the key
	_ = client.Set("mykey", "myvalue")

	// Now it should exist
	exists, err = client.Exists("mykey")
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if !exists {
		t.Error("Expected key to exist")
	}
}

func TestGetNonExistentKey(t *testing.T) {
	mock := NewMockRedditAPI()
	client := NewWithAPI(mock, "testsubreddit")

	_, err := client.Get("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent key")
	}

	if _, ok := err.(*KeyNotFoundError); !ok {
		t.Errorf("Expected KeyNotFoundError, got %T: %v", err, err)
	}
}

func TestAppendToNonExistentKey(t *testing.T) {
	mock := NewMockRedditAPI()
	client := NewWithAPI(mock, "testsubreddit")

	err := client.Append("nonexistent", "value", nil)
	if err == nil {
		t.Error("Expected error for appending to non-existent key")
	}

	if _, ok := err.(*KeyNotFoundError); !ok {
		t.Errorf("Expected KeyNotFoundError, got %T: %v", err, err)
	}
}

func TestDeleteNonExistentKey(t *testing.T) {
	mock := NewMockRedditAPI()
	client := NewWithAPI(mock, "testsubreddit")

	err := client.Delete("nonexistent")
	if err == nil {
		t.Error("Expected error for deleting non-existent key")
	}

	if _, ok := err.(*KeyNotFoundError); !ok {
		t.Errorf("Expected KeyNotFoundError, got %T: %v", err, err)
	}
}
