package file

type StaticFile struct {
	Filename     string `db:"file_path"`
	FileAsString string `db:"file_data"`
	File         []byte
}
