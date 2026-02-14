package cmd

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	dcrypt "github.com/OpenListTeam/OpenList/v4/drivers/crypt"
	"github.com/OpenListTeam/OpenList/v4/internal/db"
	"github.com/OpenListTeam/OpenList/v4/internal/fs"
	mmm "github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/OpenListTeam/OpenList/v4/internal/op"
	rcCrypt "github.com/rclone/rclone/backend/crypt"
	"github.com/rclone/rclone/fs/config/configmap"
	"github.com/rclone/rclone/fs/config/obscure"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// encryption and decryption command format for Crypt driver
// 加解密命令
// 1. 此命令可以用加密驱动的方式加密或解密本地文件/文件夹
// 2. 如果设置了--remote-names，则可以对网盘里面加密文件夹里面的文件/夹名进行加解密转换
//		比如，挂载了网盘A，并挂载了加密存储C，映射了A中的文件夹D，文件名和文件夹名都设置成了加密
//		后面觉得，在网盘A中直接浏览D中的文件很不方便，通过此命令可以直接将D中的文件名重命名为加密前的文件名而不用加密解密文件内容
//		注意：这个参数会直接修改网盘里的文件名，如果设置了cmount参数，会取消加密存储C的“文件/名加密”配置
//

type options struct {
	op          string //decrypt or encrypt
	src         string //source dir or file
	dst         string //out destination
	remoteNames bool   //指定此项，会直接修改（加密/解密）网盘里面的加密文件夹里面的文件/夹名，但不会修改文件内容（仍然是加密的）

	cmount             string //加密驱动挂载名称，如果提供了此项，会忽略下面的配置项，直接从数据库读取配置
	pwd                string //de/encrypt password
	salt               string
	filenameEncryption string //reference drivers\crypt\meta.go Addtion
	dirnameEncryption  string
	filenameEncode     string
	suffix             string
}

var opt options
var cipher *rcCrypt.Cipher

// CryptCmd represents the crypt command
var CryptCmd = &cobra.Command{
	Use:     "crypt",
	Short:   "Encrypt or decrypt local file or dir",
	Example: `openlist crypt  -s ./src/encrypt/ --op=de --pwd=123456 --salt=345678`,
	Run: func(cmd *cobra.Command, args []string) {
		opt.validate()
		opt.cryptFileDir()

	},
}

func init() {
	RootCmd.AddCommand(CryptCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	CryptCmd.Flags().StringVarP(&opt.src, "src", "s", "", "src file or dir to encrypt/decrypt")
	CryptCmd.Flags().StringVarP(&opt.dst, "dst", "d", "", "dst dir to output,if not set,output to src dir")
	CryptCmd.Flags().StringVar(&opt.op, "op", "", "de or en which stands for decrypt or encrypt")
	CryptCmd.Flags().BoolVar(&opt.remoteNames, "remote-names", false, "only rename names of remote files")

	CryptCmd.Flags().StringVar(&opt.cmount, "cmount", "", "the name of mounted crypt in db")
	CryptCmd.Flags().StringVar(&opt.pwd, "pwd", "", "password used to encrypt/decrypt,if not contain ___Obfuscated___ prefix,will be obfuscated before used")
	CryptCmd.Flags().StringVar(&opt.salt, "salt", "", "salt used to encrypt/decrypt,if not contain ___Obfuscated___ prefix,will be obfuscated before used")
	CryptCmd.Flags().StringVar(&opt.filenameEncryption, "filename-encrypt", "off", "filename encryption mode: off,standard,obfuscate")
	CryptCmd.Flags().StringVar(&opt.dirnameEncryption, "dirname-encrypt", "false", "is dirname encryption enabled:true,false")
	CryptCmd.Flags().StringVar(&opt.filenameEncode, "filename-encode", "base64", "filename encoding mode: base64,base32,base32768")
	CryptCmd.Flags().StringVar(&opt.suffix, "suffix", ".bin", "suffix for encrypted file,default is .bin")
}

func (o *options) validate() {
	if o.src == "" {
		log.Fatal("src can not be empty")
	}
	if o.op != "encrypt" && o.op != "decrypt" && o.op != "en" && o.op != "de" {
		log.Fatal("op must be encrypt or decrypt")
	}
	if o.cmount == "" && o.filenameEncryption != "off" && o.filenameEncryption != "standard" && o.filenameEncryption != "obfuscate" {
		log.Fatal("filename_encryption must be off,standard,obfuscate")
	}
	if o.cmount == "" && o.filenameEncode != "base64" && o.filenameEncode != "base32" && o.filenameEncode != "base32768" {
		log.Fatal("filename_encode must be base64,base32,base32768")
	}

}

func (o *options) cryptFileDir() {
	src, _ := filepath.Abs(o.src)
	log.Printf("src abs is %v", src)

	if o.remoteNames {

		Init()
		if o.cmount != "" {
			ms, err := db.GetStorageByMountPath(o.cmount)
			if err != nil {
				log.Fatalf("can't find mount path %v, err: %v\n", o.cmount, err)
			}
			var ad dcrypt.Addition
			err = json.Unmarshal([]byte(ms.Addition), &ad)
			if err != nil {
				log.Fatalf("unmarshal crypt addition err:%v \n", err)
			}
			createCipher(&ad)
		} else {
			createCipher(nil)
		}
		cryptNames(o.src)
		return
	}

	fileInfo, err := os.Stat(src)
	if err != nil {
		log.Fatalf("reading file/dir %v failed,err:%v", src, err)

	}
	createCipher(nil)
	dst := ""
	//check and create dst dir
	if o.dst != "" {
		dst, _ = filepath.Abs(o.dst)
		checkCreateDir(dst)
	}

	// src is file
	if !fileInfo.IsDir() { //file
		if dst == "" {
			dst = filepath.Dir(src)
		}
		o.cryptFile(cipher, src, dst)
		return
	}

	// src is dir
	if dst == "" {
		//if src is dir and not set dst dir ,create ${src}_crypt dir as dst dir
		dst = path.Join(filepath.Dir(src), fileInfo.Name()+"_crypt")
	}
	log.Printf("dst : %v\n", dst)

	dirnameMap := make(map[string]string)
	pathSeparator := string(os.PathSeparator)

	filepath.Walk(src, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("get file %v info failed, err:%v\n", p, err)
		}
		if p == src {
			return nil
		}
		log.Printf("current path %v\n", p)

		// relative path
		rp := strings.ReplaceAll(p, src, "")
		log.Printf("relative path %v\n", rp)

		rpds := strings.Split(rp, pathSeparator)

		if info.IsDir() {
			// absolute dst dir for current path
			dd := ""

			if o.dirnameEncryption == "true" {
				if o.op == "encrypt" || o.op == "en" {
					for i := range rpds {
						oname := rpds[i]
						if _, ok := dirnameMap[rpds[i]]; ok {
							rpds[i] = dirnameMap[rpds[i]]
						} else {
							rpds[i] = cipher.EncryptDirName(rpds[i])
							dirnameMap[oname] = rpds[i]
						}
					}
					dd = path.Join(dst, strings.Join(rpds, pathSeparator))
				} else {
					for i := range rpds {
						oname := rpds[i]
						if _, ok := dirnameMap[rpds[i]]; ok {
							rpds[i] = dirnameMap[rpds[i]]
						} else {
							dnn, err := cipher.DecryptDirName(rpds[i])
							if err != nil {
								log.Fatalf("decrypt dir name %v failed,err:%v\n", rpds[i], err)
							}
							rpds[i] = dnn
							dirnameMap[oname] = dnn
						}

					}
					dd = path.Join(dst, strings.Join(rpds, pathSeparator))
				}

			} else {
				dd = path.Join(dst, rp)
			}

			log.Printf("create output dir %v", dd)
			checkCreateDir(dd)
			return nil
		}

		// file dst dir
		fdd := dst

		if o.dirnameEncryption == "true" {
			for i := range rpds {
				if i == len(rpds)-1 {
					break
				}
				fdd = path.Join(fdd, dirnameMap[rpds[i]])
			}

		} else {
			fdd = path.Join(fdd, strings.Join(rpds[:len(rpds)-1], pathSeparator))
		}

		log.Printf("file output dir %v", fdd)
		o.cryptFile(cipher, p, fdd)
		return nil
	})

}

func (o *options) cryptFile(cipher *rcCrypt.Cipher, src string, dst string) {
	fileInfo, err := os.Stat(src)
	if err != nil {
		log.Fatalf("get file %v  info failed,err:%v", src, err)

	}
	fd, err := os.OpenFile(src, os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("open file %v failed,err:%v", src, err)

	}
	defer fd.Close()

	var cryptSrcReader io.Reader
	var outFile string
	if o.op == "encrypt" || o.op == "en" {
		filename := fileInfo.Name()
		if o.filenameEncryption != "off" {
			filename = cipher.EncryptFileName(fileInfo.Name())
			log.Printf("encrypt file name %v to %v\n", fileInfo.Name(), filename)
		} else {
			filename = fileInfo.Name() + o.suffix
		}
		cryptSrcReader, err = cipher.EncryptData(fd)
		if err != nil {
			log.Fatalf("encrypt file %v failed,err:%v", src, err)

		}
		outFile = path.Join(dst, filename)
	} else {
		filename := fileInfo.Name()
		if o.filenameEncryption != "off" {
			filename, err = cipher.DecryptFileName(filename)
			if err != nil {
				log.Fatalf("decrypt file name %v failed,err:%v", src, err)
			}
			log.Printf("decrypt file name %v to %v \n", fileInfo.Name(), filename)
		} else {
			filename = strings.TrimSuffix(filename, o.suffix)
		}

		cryptSrcReader, err = cipher.DecryptData(fd)
		if err != nil {
			log.Fatalf("decrypt file %v failed,err:%v", src, err)

		}
		outFile = path.Join(dst, filename)
	}
	//write new file
	wr, err := os.OpenFile(outFile, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Fatalf("create file %v failed,err:%v", outFile, err)

	}
	defer wr.Close()

	_, err = io.Copy(wr, cryptSrcReader)
	if err != nil {
		log.Fatalf("write file %v failed,err:%v", outFile, err)
	}

}

// check dir exist ,if not ,create
func checkCreateDir(dir string) {
	_, err := os.Stat(dir)

	if os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatalf("create dir %v failed,err:%v", dir, err)
		}
		return
	} else if err != nil {
		log.Fatalf("read dir %v err: %v", dir, err)
	}

}

func updateObfusParm(str string) string {
	obfuscatedPrefix := "___Obfuscated___"
	if !strings.HasPrefix(str, obfuscatedPrefix) {
		str, err := obscure.Obscure(str)
		if err != nil {
			log.Fatalf("update obfuscated parameter failed,err:%v", str)
		}
	} else {
		str, _ = strings.CutPrefix(str, obfuscatedPrefix)
	}
	return str
}

func createCipher(ad *dcrypt.Addition) {
	if ad != nil {
		opt.pwd = ad.Password
		opt.salt = ad.Salt
		opt.filenameEncryption = ad.FileNameEnc
		opt.dirnameEncryption = ad.DirNameEnc
		opt.filenameEncode = ad.FileNameEncoding
		opt.suffix = ad.EncryptedSuffix

	}
	pwd := updateObfusParm(opt.pwd)
	salt := updateObfusParm(opt.salt)
	config := configmap.Simple{
		"password":                  pwd,
		"password2":                 salt,
		"filename_encryption":       opt.filenameEncryption,
		"directory_name_encryption": opt.dirnameEncryption,
		"filename_encoding":         opt.filenameEncode,
		"suffix":                    opt.suffix,
		"pass_bad_blocks":           "",
	}

	c, err := rcCrypt.NewCipher(config)
	if err != nil {
		log.Fatalf("create cipher err: %v", err)
	}
	cipher = c

}

func cryptNames(path string) {
	log.Printf("deal path： %s\n", path)
	objs, err := fs.List(context.TODO(), path, &fs.ListArgs{
		Refresh:            true,
		WithStorageDetails: false,
	})
	if err != nil {
		log.Printf("get list err dir name err:%v\n", err)
		return

	}
	log.Printf("get list success,path: %v\n", path)
	var renameObjects []mmm.RenameObj

	for _, obj := range objs {
		isdir := obj.IsDir()
		name := obj.GetName()
		if isdir {
			time.Sleep(time.Second) //防止请求频率过快被限制
			cryptNames(filepath.Join(path, name))
		}
		var newname string
		if opt.op == "en" { //加密
			newname = strings.TrimSuffix(name, opt.suffix)
			if isdir {
				newname = cipher.EncryptDirName(newname)
			} else {
				newname = cipher.EncryptFileName(newname)
			}

		} else { //解密
			if isdir {
				tn, err := cipher.DecryptDirName(name)

				if err != nil {
					log.Fatalf("decrypt dir name %v err: %v", name, err)
				}
				newname = tn
			} else {
				tn, err := cipher.DecryptFileName(name)
				if err != nil {
					log.Fatalf("decrypt filename %v err: %v", name, err)
				}
				newname = tn + opt.suffix //文件名拼接后缀
			}

		}

		renameObjects = append(renameObjects, mmm.RenameObj{
			ID:      obj.GetID(),
			SrcName: obj.GetName(),
			NewName: newname,
		})

	}
	if len(renameObjects) == 0 {
		return
	}
	//执行批量重命名
	storage, actualPath, err := op.GetStorageAndActualPath(path)
	err = fs.BatchRename(context.Background(), storage, actualPath, renameObjects)
	if err != nil {
		log.Fatalf("batchrename error : %v\n", err)
	}
	time.Sleep(time.Second) //防止请求频率过快被限制

}
