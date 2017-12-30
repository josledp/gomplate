// +build !windows

package main

// func TestReadInput(t *testing.T) {
// 	actual, err := readInputs("foo", nil)
// 	assert.Nil(t, err)
// 	assert.Equal(t, &input{"<arg>", "foo"}, actual[0])

// 	// stdin is "" because during tests it's given /dev/null
// 	actual, err = readInputs("", []string{"-"})
// 	assert.Nil(t, err)
// 	assert.Equal(t, &input{"-", ""}, actual[0])

// 	actual, err = readInputs("", []string{"template_test.go"})
// 	assert.Nil(t, err)
// 	thisFile, _ := os.Open("template_test.go")
// 	expected, _ := ioutil.ReadAll(thisFile)
// 	assert.Equal(t, &input{"template_test.go", string(expected)}, actual[0])
// }

// func TestInputDir(t *testing.T) {
// 	outDir, err := ioutil.TempDir(filepath.Join("test", "files", "input-dir"), "out-temp-")
// 	assert.Nil(t, err)
// 	defer (func() {
// 		if cerr := os.RemoveAll(outDir); cerr != nil {
// 			log.Fatalf("Error while removing temporary directory %s : %v", outDir, cerr)
// 		}
// 	})()

// 	src, err := data.ParseSource("config=test/files/input-dir/config.yml")
// 	assert.Nil(t, err)

// 	d := &data.Data{
// 		Sources: map[string]*data.Source{"config": src},
// 	}
// 	gomplate := NewGomplate(d, "{{", "}}")
// 	err = processInputDir(filepath.Join("test", "files", "input-dir", "in"), outDir, []string{"**/*.exclude.txt"}, gomplate)
// 	assert.Nil(t, err)

// 	top, err := ioutil.ReadFile(filepath.Join(outDir, "top.txt"))
// 	assert.Nil(t, err)
// 	assert.Equal(t, "eins", string(top))

// 	inner, err := ioutil.ReadFile(filepath.Join(outDir, "inner/nested.txt"))
// 	assert.Nil(t, err)
// 	assert.Equal(t, "zwei", string(inner))

// 	// excluded file should not exist in out dir
// 	_, err = ioutil.ReadFile(filepath.Join(outDir, "inner/exclude.txt"))
// 	assert.NotEmpty(t, err)
// }
