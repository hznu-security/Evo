package fs

//type frontendFS struct {
//	frontendFS http.FileSystem
//}
//
//func FS() *frontendFS {
//	return &frontendFS{
//		frontendFS: New(),
//	}
//}
//
//func (f *frontendFS) Open(name string) (http.File, error) {
//	return f.frontendFS.Open(name)
//}
//
//func (f *frontendFS) Exists(prefix string, filepath string) bool {
//	if _, err := f.frontendFS.Open(filepath); err != nil {
//		return false
//	}
//	return true
//}
