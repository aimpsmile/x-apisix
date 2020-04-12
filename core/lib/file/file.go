package file

import (
	"os"
	"path/filepath"
	"runtime"

	mlog "github.com/micro/go-micro/v2/logger"
)

//	获取运行脚本的目录
func GetRunDir(skip int) string {
	_, file, _, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	dir, err := filepath.Abs(filepath.Dir(file))
	if err != nil {
		mlog.Errorf("core.lib.file.dir.err %v", err)
		return ""
	}
	return dir
}

//	true 文件存在，err不为nil
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

//	覆盖写入
func WriteToFile(fileName, content string, perm os.FileMode) error {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, perm)
	if err != nil {
		mlog.Errorf("core.lib.file.write.err %v", err)
		return err
	} else {
		_, err = f.Write([]byte(content))
		if e := f.Close(); e != nil {
			mlog.Errorf("core.lib.file.close.err %v", e)
		}
		return err
	}
}

//	追加写入
func AppendToFile(fileName, content string, perm os.FileMode) error {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, perm)
	if err != nil {
		mlog.Errorf("core.lib.file.append.err %v", err)
		return err
	} else {
		_, err = f.Write([]byte(content))
		if e := f.Close(); e != nil {
			mlog.Errorf("core.lib.file.close.err %v", e)
		}
		return err
	}
}
