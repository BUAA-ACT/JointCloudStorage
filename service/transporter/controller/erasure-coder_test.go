package controller

import (
	"os"
	"testing"
)

func TestEncodeAndDecode(t *testing.T) {
	shards := []string{
		"../test/tmp/testEncode.txt.0", "../test/tmp/testEncode.txt.1", "../test/tmp/testEncode.txt.2",
	}
	filePath := "../test/tmp/test.txt"
	fi, _ := os.Stat(filePath)
	fileSize := fi.Size()
	t.Logf("file size: %v", fileSize)
	t.Run("TestEncode", func(t *testing.T) {
		err := Encode(filePath, shards, 2, 1)
		if err != nil {
			t.Errorf("test encode err: %v", err)
		}
	})
	t.Run("testDecode", func(t *testing.T) {
		shards[2] = "../test/tmp/testEncode.txt.2.fail"
		err := Decode("../test/tmp/testEncode.txt", fileSize, shards, 2, 1)
		if err != nil {
			t.Fatalf("test decode err: %v", err)
		}
		fi, _ := os.Stat("../test/tmp/testEncode.txt")
		t.Logf("file size: %v", fi.Size())
	})

}
