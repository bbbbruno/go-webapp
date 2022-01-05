package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/objx"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	client := new(client)
	_, err := authAvatar.GetAvatarURL(client)
	if err != ErrNoAvatarURL {
		t.Error("値が存在しない場合、AuthAvatar.GetAvatarURLは" + "ErrNoAvatarURLを返すべきです")
	}

	testUrl := "http://url-to-avatar/"
	client.userData = objx.Map{"avatar_url": testUrl}
	url, err := authAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("値が存在する場合、AuthAvatar.GetAvatarURLは " + "エラーを返すべきではありません")
	} else {
		if url != testUrl {
			t.Error("AuthAvatar.GetAvatarURLは正しいURLを返すべきです")
		}
	}
}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	client := new(client)
	client.userData = objx.Map{"userid": "0bc83cb571cd1c50ba6f3e8a78ef1346"}
	url, err := gravatarAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("GravatarAvatar.GetAvatarURLはエラーを返すべきではありません")
	}
	if url != "//www.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78ef1346" {
		t.Errorf("GravatarAvatar.GetAvatarURLが%sという誤った値を返しました", url)
	}
}

func TestFileSystemAvatar(t *testing.T) {
	filename := filepath.Join("avatars", "abc.jpg")
	ioutil.WriteFile(filename, []byte{}, 0600)
	defer func() { os.Remove(filename) }()

	var fileSystemAvatar FileSystemAvatar
	client := new(client)
	client.userData = objx.Map{"userid": "abc"}
	url, err := fileSystemAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("ileSystemAvatar.GetAvatarURLはエラーを返すべきではありません")
	}
	if url != "/avatars/abc.jpg" {
		t.Errorf("ileSystemAvatar.GetAvatarURLが%sという誤った値を返しました", url)
	}
}
