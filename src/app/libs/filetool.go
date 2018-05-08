package libs

import (
	"path"
	"strings"
	"github.com/satori/go.uuid"
	"fmt"
	"os"
	"io"
	"archive/tar"
	"compress/gzip"
	
)
/*此结构没有写好,需要重构*/

type FileTool struct {
	Url string
}

// 获取文件名带后缀
func (this *FileTool) FileNameWithExt() string  {
	if this.Url == "" {
		return ""
	}
	return path.Base(this.Url)
}

// 获取文件后缀
func (this *FileTool) Ext() string {
	fileName := this.FileNameWithExt()
	if fileName == "" {
		return ""
	}
	return path.Ext(fileName)
}

// 获取文件后缀
func (this *FileTool) FileName() string {
	fileName := this.FileNameWithExt()
	if fileName == "" {
		return ""
	}
	return strings.TrimSuffix(fileName, this.Ext())
}

// 生成UUID的文件名
func (this *FileTool) CreateUuidFile() string  {
	ext := this.Ext()
	if ext == "" {
		return ""
	}
	return fmt.Sprintf("%s%s",uuid.NewV4().String(), ext)
}

func (this *FileTool) GenerateUuidStr() string {
	return uuid.NewV4().String()
}

// 判断文件类型是否为想要类型
func (this *FileTool) CheckFileExt(exts []string) bool  {
	if len(exts) == 0 {
		return true
	}
	ext := strings.TrimLeft(this.Ext(), ".")

	for _, value := range exts {
		if ext == value {
			return true
		}
	}
	return false
}

//判断文件是否存在
func (this *FileTool) IsExist() bool {
	if (this.Url == "") {
		return false
	}
	_, err := os.Stat(this.Url)
	return err == nil || os.IsExist(err)
}

func (this *FileTool) Ungzip(sourcefile, destDirPath string) error {
	os.Mkdir(destDirPath, os.ModePerm)

    sf, err := os.Open(sourcefile)
    if err != nil {
		return  nil
	}
    defer sf.Close()

    // Gzip reader
    gr, err := gzip.NewReader(sf)
	
	if err != nil {
		return err
	}

    // Tar reader
    tr := tar.NewReader(gr)

    for {
        hdr, err := tr.Next()
        if err == io.EOF {
            break
        }
        // Check if it is diretory or file
        if hdr.Typeflag != tar.TypeDir {
            // Get files from archive
            // Create diretory before create file
            os.MkdirAll(destDirPath+"/"+path.Dir(hdr.Name), os.ModePerm)
            fw, _ := os.Create(destDirPath + "/" + hdr.Name)
            _, err = io.Copy(fw, tr)
			if err != nil {
				return err
			}
        }
    }
	return nil
}

func (this *FileTool) DeleteDir(dir string) {
	if dir == "" || dir == " " {
		return 
	}
	os.RemoveAll(dir)
}






