package model

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

var __1_data_down_sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\xce\x4b\x0e\xc3\x20\x0c\x04\xd0\x7d\x4e\x91\x7b\xe4\x30\x16\x9f\x09\xb5\x04\xa6\xb2\x4d\xcf\xdf\x45\xbb\x09\xb0\x9d\x67\x6b\x26\x6b\x7f\x9f\x1e\x62\xc5\x89\x16\x91\x8d\x6e\x46\xcd\x76\x1d\x8b\x3c\xa2\x06\x71\xee\x42\xc3\xa0\x7b\xd1\x5e\x31\x8b\x59\x28\x53\x58\x06\xd7\x6c\xf4\xe9\x9c\x40\xe6\xc1\xc7\xd6\x37\x92\x5e\x41\x04\x75\xee\x68\x71\x5e\xb4\x6e\x5c\xb7\xfd\x6a\xae\xe3\x51\xd0\xe5\xe6\x42\x0a\x73\xe5\xe4\xa4\x63\xfe\xfa\x5f\x84\xdc\x58\x76\x72\x7d\x03\x00\x00\xff\xff\xfb\x07\x86\x29\x5f\x01\x00\x00")

func _1_data_down_sql() ([]byte, error) {
	return bindata_read(
		__1_data_down_sql,
		"1_data.down.sql",
	)
}

var __1_data_up_sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x59\xd1\x96\xa4\x28\x0c\x7d\xaf\xaf\xe0\xb1\x7b\xcf\x9e\xf9\x81\xed\xfd\x16\x0f\x6a\xac\xca\x36\x82\x07\xb0\xaa\xeb\xef\xf7\x28\x88\x80\x04\x9d\x99\x27\xcd\xbd\x09\x98\x1b\x48\xaa\xa7\xd3\xc0\x2d\x30\xcb\x5b\x01\xac\x53\x72\xc0\xfb\xed\xe3\xc6\xfc\x63\x73\x9f\x51\xf4\x4d\x8f\xa6\x53\xba\x6f\xb0\x67\x16\x7e\x2c\x9b\x34\x8e\x5c\xbf\xd9\x37\xbc\xff\xde\xb9\x93\x86\x01\x7f\x56\x46\x64\xed\x94\x50\xba\x01\xad\x95\x66\x2d\xde\x51\x1e\xc1\x17\xd7\x92\xc2\x50\x0e\xea\x88\xf1\xd9\x2a\x0d\xa3\x7a\x42\xc3\x3b\x8b\x4f\x60\xad\x52\x02\xb8\x2c\x73\x2c\x8e\x70\x0c\xa2\xc1\x58\x8d\x9d\xcd\x42\xdc\x3e\xff\xb9\xdd\x0a\x69\x69\x78\x3f\xa2\x34\x71\x76\x9c\x65\x49\x4b\x8b\x77\x03\x1a\xb9\x20\x72\xe3\x99\xc5\x74\x1e\x69\xf6\x3d\x15\xf6\xeb\xc1\x09\xf4\x88\xc6\xa0\x92\xcd\xc8\xcd\x37\xc9\xd3\x4a\xc0\xf9\x52\xb3\x01\x9d\xb3\xc8\x0c\x84\x8c\xe9\x59\x40\x92\x89\x14\xb9\x92\x91\xcc\xe3\x2c\x33\x19\x1d\x9e\x20\x6d\x45\x52\xc7\x32\x20\xa0\xb3\x4a\x17\xbf\xf2\x82\x5b\xf7\xe0\x52\x82\xf8\x03\xcf\x75\xc1\x89\x5b\x0b\x5a\xfe\xee\x6a\xd7\xdc\x96\xa2\x55\x85\x63\x53\xa4\x35\x52\x59\x1c\xde\x85\x53\x52\xa3\xff\xfe\xf7\x7b\xf7\x93\xf3\x96\x70\x4b\x65\x7a\x2c\xc0\xb5\x3c\x5c\xc5\xb9\xc7\x6a\x89\x79\x4a\x61\xdb\x1e\x91\x7c\x84\xdc\x86\x23\xbf\x03\x6b\xdf\x16\x78\x64\x35\x93\xe0\xe6\x71\x30\xbb\xcd\xf5\x6c\xf9\x52\x63\xf9\x38\xb1\x1e\x06\x3e\x0b\xcb\xa4\x7a\x7d\x7c\xc6\x9b\x00\x01\x09\x33\xfe\xb8\x27\xc2\x6b\x23\x0a\x6e\x2c\xe3\xe6\xe6\xea\x81\xf1\x5f\x7f\xb1\x41\xab\xd1\xc3\x8c\x33\x01\x83\x65\xff\x29\x94\x9b\xa9\x65\x4a\xde\xd8\xfa\x8f\xff\x3a\x7e\xf3\xbf\xac\x2d\x58\xb9\xec\x73\x1f\xec\xd9\xd7\xce\xc5\xfe\xf6\x7a\x80\x86\xd8\xc2\xd0\x30\x39\x0b\x91\xcb\xb2\x88\xe7\x54\x59\x9f\xaa\xa2\x38\x06\x79\xce\x1d\x4c\x02\x89\x62\xce\xb4\x36\x88\xa8\xcc\x9c\x75\x52\x06\xb3\x93\xe1\x81\x70\x6f\x9a\x03\x76\x41\x4e\xbf\xbf\x53\x35\x1d\x8f\x12\x73\x45\x13\x2d\x9d\x25\x91\xf2\x90\x8a\x45\xc9\x83\x31\x16\x32\x64\xff\x2b\x30\x23\x19\x03\x4a\xa8\xb8\x5c\x56\x4e\xc5\xf5\xa9\xaa\xa2\x63\x14\x64\x72\x40\x22\xd3\xce\xd5\x38\xa2\xe4\x56\xe9\x0c\x6b\x95\x8d\xef\x24\x67\x84\x91\xa3\xc8\x88\xfc\xc9\x2d\xd7\xfb\x31\x74\x56\xa1\x3a\x2e\xf2\xf5\x9e\xa0\x71\x40\xe8\x8f\x81\x2f\xa8\xec\xb7\x7c\xaa\xb2\x5f\x9f\x50\x79\x45\x13\x95\x9d\x25\x51\xf9\x90\xc9\x45\xe5\x83\x31\x56\x39\xa8\xf3\x15\x98\x91\xca\x01\x25\x54\x1e\x61\x6c\x37\x9d\xfd\x73\x55\xe9\x8d\x43\x9e\xd8\x8d\x40\x35\xd7\x0d\x97\xd8\x7d\x1f\x8c\x4b\x56\xa0\x6f\xb8\xdd\x93\x1c\xc3\x17\xa4\xda\xa8\xe7\x62\x6d\x4c\x4a\x2e\x8f\x27\x82\x6d\xb6\x44\x32\x32\x25\x8b\x74\x24\x18\x4b\x48\xe5\x2c\x0e\x90\x63\x25\x7f\x5f\x04\xfb\x6b\x28\x83\x88\x41\x14\x82\xef\xea\x7e\x7e\xf3\x2f\xf5\x91\x6d\x23\xd1\x53\xda\xc6\xa8\x61\xc9\xe5\x10\xac\xf9\xb0\xbb\xd9\x5b\xb4\x7a\xd9\x74\x01\x9a\xb8\x06\x69\xab\x6b\x15\xda\xc0\xbe\x0f\x33\xbc\x92\x51\x28\x6c\x45\x4d\xd8\x1d\x63\xad\x82\x08\x1c\xd1\x96\xa2\x5d\x28\xd5\x3d\x3b\xa7\xb5\x1a\xa8\x54\xb1\x6e\x84\xa4\x5a\x83\x31\x29\xd7\x92\x28\x4b\xa1\x95\xec\x71\x91\xc5\x35\xf1\x15\xf3\xa3\x32\x8b\x39\x44\x9d\xf9\xe9\xe1\xa9\xb0\x83\xc6\x58\x6e\xe7\x64\x80\x8b\xed\x57\xa6\xb9\x84\xbf\x4e\x98\xe7\xc3\x57\xe2\x43\x16\x6f\x89\x5c\x19\x7d\x4b\x74\xea\x0e\x2c\x86\xbe\x3c\x3b\x26\x6e\x97\x07\xc9\xc4\xab\x3e\x55\x26\xdc\xc2\x88\x99\xe2\xa5\x79\xb3\x9e\xe1\x68\x04\xad\x13\x0b\x53\x69\x35\xbf\x54\xe0\xda\xdd\x49\xd4\xdd\x17\x11\xea\x38\x07\xe7\x9e\x64\xa3\x35\x86\xdf\x61\xeb\xb4\xee\xe5\xa4\xd5\x7a\x52\xa5\xd7\x7a\x46\xa5\x2e\x03\x87\x6e\xc8\x9e\x50\xc3\xb2\xeb\x38\xd8\x5f\xd0\x3e\x94\xfa\xae\x6f\x4e\x49\xbb\xfc\x2e\x3f\x02\x57\xda\xf9\xb6\xb9\x0b\xfd\xdc\x53\xe9\x86\xee\x08\x59\x47\xf7\xc6\xac\xa5\x53\xa9\x77\x2d\x99\x42\xd3\xa6\x4c\x8b\x93\x44\x29\xe0\xc5\x38\xe5\xe9\x80\x00\x8b\x11\x28\xe7\x33\xbf\x30\x57\x84\xf7\x68\xb0\xd8\x39\x44\xe5\x2f\xa3\x87\xff\x95\xee\x1e\xab\x55\xef\x29\x3e\xae\xa7\xa2\xb4\x4c\xc3\x00\x1a\x64\x07\x26\x88\xf6\x11\xad\xfe\x19\x39\xef\x05\xb5\xd5\xdc\x06\x2c\x65\x9c\xd9\xf2\x1f\x8c\xde\xdc\xc3\xf2\x03\x65\x5a\xa7\x85\x3c\x0a\x5a\x71\x08\x33\x6b\x91\x9b\xf8\x6c\x1f\x4a\xa7\x03\x4e\x0a\x61\xa7\x64\xc5\x73\xd2\xea\xe7\x7d\x46\x2a\x20\x93\x56\x4f\xec\xa1\xbc\x74\x00\x0b\x8e\xf6\x31\x8f\xad\xe4\x28\x9a\x17\xf6\xf6\x71\xcc\xca\x4e\x78\x00\xde\x1f\xb6\xc6\xa8\xc6\x77\x9f\x56\xa0\xac\x7f\x70\xa1\x96\x77\x20\xb5\xb4\x43\xc9\x98\xe4\x92\x4b\x3a\x14\xb5\xa4\x03\xa9\x25\x1d\x4a\xc6\x24\x97\x1c\x94\xb2\xa0\x9b\xf5\xef\xe6\x65\x88\x52\xdd\xc3\xf5\xd2\xb8\x70\xb3\x86\x32\xaf\xdc\xab\xf1\xf9\x6d\x06\x04\x91\x1e\x63\x67\xb9\x72\x9a\x3d\x73\x7d\x2b\x1f\x69\x47\xfc\x08\x17\xc4\xe7\xd1\xbb\x54\xcc\x1e\x7a\x72\x31\x53\x18\x4a\x81\x32\xf9\xaf\x80\x14\xbf\x9e\x2b\xef\x50\x4c\xd9\xff\x01\x00\x00\xff\xff\xc7\xdf\x31\xdf\x32\x19\x00\x00")

func _1_data_up_sql() ([]byte, error) {
	return bindata_read(
		__1_data_up_sql,
		"1_data.up.sql",
	)
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() ([]byte, error){
	"1_data.down.sql": _1_data_down_sql,
	"1_data.up.sql":   _1_data_up_sql,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func     func() ([]byte, error)
	Children map[string]*_bintree_t
}

var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"1_data.down.sql": &_bintree_t{_1_data_down_sql, map[string]*_bintree_t{}},
	"1_data.up.sql":   &_bintree_t{_1_data_up_sql, map[string]*_bintree_t{}},
}}
