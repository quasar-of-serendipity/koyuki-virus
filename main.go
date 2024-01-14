package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/sys/windows/registry"
)

//go:embed koyuki.wav
var koyuki []byte

const NAMESPACE = "NIHAHAHAHA"
const NAMESPACE_SHORT = "NIHAH0"
const NAMESPACE_CURRENT = ".Current"

const SCHEMES_KEY_BASE = "AppEvents\\Schemes"
const NAMES_KEY_BASE = "AppEvents\\Schemes\\Names"
const APPS_KEY_BASE = "AppEvents\\Schemes\\Apps"

func makePath(components ...string) string {
	return strings.Join(components, "\\")
}

func getKoyukiPath() string {
	return os.Getenv("USERPROFILE") + "\\Documents\\koyuki.wav"
}

func createSoundProfile() error {
	key, _, err := registry.CreateKey(registry.CURRENT_USER, makePath(NAMES_KEY_BASE, NAMESPACE_SHORT), registry.ALL_ACCESS)
	if err != nil {
		return err
	}

	return key.SetStringValue("", NAMESPACE)
}

func createSoundAssociations(namespace string) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, APPS_KEY_BASE, registry.ALL_ACCESS)
	if err != nil {
		return err
	}

	names, err := key.ReadSubKeyNames(0)
	if err != nil {
		return err
	}

	for _, name := range names {
		err = createSoundAssociationsInSubkey(namespace, name)
		if err != nil {
			return err
		}
	}

	return nil
}

func createSoundAssociationsInSubkey(namespace, subkey string) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, makePath(APPS_KEY_BASE, subkey), registry.ALL_ACCESS)
	if err != nil {
		return err
	}

	names, err := key.ReadSubKeyNames(0)
	if err != nil {
		return err
	}

	for _, name := range names {
		err = makeKoyukiKey(namespace, subkey, name)
		if err != nil {
			return err
		}
	}

	return nil
}

func makeKoyukiKey(namespace, subkey, name string) error {
	path := makePath(APPS_KEY_BASE, subkey, name, namespace)
	key, _, err := registry.CreateKey(registry.CURRENT_USER, path, registry.ALL_ACCESS)
	if err != nil {
		return err
	}

	return key.SetStringValue("", getKoyukiPath())
}

func setDefaultSoundScheme() error {
	key, err := registry.OpenKey(registry.CURRENT_USER, SCHEMES_KEY_BASE, registry.ALL_ACCESS)
	if err != nil {
		return err
	}

	return key.SetStringValue("", NAMESPACE_SHORT)
}

func printKoyuki() {
	fmt.Println("\n\t\t    NIHAHAHAHA!")
	fmt.Println()

	fmt.Print(
		"....,.........................*............*......\n" +
			"...,............*............../............,,....\n" +
			".,,.............*...............,.............*...\n" +
			",,.............,*...............(..............,..\n" +
			",,.....,......,/*...............(,..............,/\n" +
			",.....,......,*..,..............,.(,.............*\n" +
			",....,......./   /.....,......../ *.*,.....*......\n" +
			",..,,......,,  ,&%%*..,*.....%%%,     *,,...,,....\n" +
			"..,,.....,,       ....,*........        ,*,,,.,,..\n" +
			",,,..../.          ...,.,......,.,  ,#&@@@&&(,**(,\n" +
			",***.%@@&@&&&@@/    .,. ,.....,  ,&&###,  / .&&&&&\n" +
			",.(&&&  %###(#  .%.  ,  ,....,..* @##/.###%   %%%.\n" +
			",.(%%.  %##, .,(@       /....     &(, ../(&   &   \n" +
			"*.  ./  *(/*,.//&       *..*       .**,*/%.       \n" +
			",,.       &/*,*#.       ,,.                .     .\n" +
			".**.     ..                           ......... ..\n" +
			".,*.  ........                         .        ,*\n" +
			"*,,,,.                                        .%&.\n" +
			"*%%###(                                    .######\n" +
			"#########/..                          .,   .(#####",
	)

	fmt.Println()
	fmt.Println()
}

func main() {
	err := os.WriteFile(getKoyukiPath(), koyuki, 0644)
	if err != nil {
		log.Fatalf("Failed to write koyuki.wav to Documents: %v", err)
	}

	err = createSoundProfile()
	if err != nil {
		log.Fatalf("Failed to create sound profile: %v", err)
	}

	err = setDefaultSoundScheme()
	if err != nil {
		log.Fatalf("Failed to set default sound scheme: %v", err)
	}

	err = createSoundAssociations(NAMESPACE_SHORT)
	if err != nil {
		log.Fatalf("Failed to create sound associations (sound scheme): %v", err)
	}

	err = createSoundAssociations(NAMESPACE_CURRENT)
	if err != nil {
		log.Fatalf("Failed to create sound associations (current scheme): %v", err)
	}

	printKoyuki()
	time.Sleep(5 * time.Second)
}
