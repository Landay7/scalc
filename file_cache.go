package main

type fileCache map[string][]int

func (c fileCache) get(filename string) ([]int, error) {
	if nums, ok := c[filename]; ok {
		return nums, nil
	}
	nums, err := readNumbersInFile(filename)
	if err != nil {
		return nil, err
	}
	c[filename] = nums
	return nums, nil

}

func newFileCache() fileCache {
	return make(fileCache)
}
