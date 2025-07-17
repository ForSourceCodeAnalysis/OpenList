package open123

// // 创建文件 V2
// func (d *Open123) create(parentFileID int64, filename string, etag string, size int64, duplicate int, containDir bool) (*UploadCreateResp, error) {
// 	var resp UploadCreateResp
// 	_, err := d.Request(UploadCreate, http.MethodPost, func(req *resty.Request) {
// 		req.SetBody(base.Json{
// 			"parentFileId": parentFileID,
// 			"filename":     filename,
// 			"etag":         strings.ToLower(etag),
// 			"size":         size,
// 			"duplicate":    duplicate,
// 			"containDir":   containDir,
// 		})
// 	}, &resp)
// 	if err != nil {
// 		return nil, err
// 	"net/http"

// 	"github.com/go-resty/resty/v2"
// )

// func (d *Open123) preup(requ *PreupReq) (*PreupResp, error) {
// 	r := &BaseResp{
// 		Code: -1,
// 		Data: &PreupResp{},
// 	}
// 	_, err := d.Request(d.qpsInstance[preupCreateAPI], http.MethodPost, func(req *resty.Request) {
// 		req.SetBody(requ).SetHeader("Content-Type", "application/json")

// 	}, r)

// 	return r.Data.(*PreupResp), err
// }

// 上传分片 V2
// func (d *Open123) Upload(ctx context.Context, file model.FileStreamer, createResp *UploadCreateResp, up driver.UpdateProgress) error {
// 	uploadDomain := createResp.Data.Servers[0]
// 	size := file.GetSize()
// 	chunkSize := createResp.Data.SliceSize

// 	ss, err := stream.NewStreamSectionReader(file, int(chunkSize), &up)
// 	if err != nil {
// 		return err
// 	}

// 	uploadNums := (size + chunkSize - 1) / chunkSize
// 	thread := min(int(uploadNums), d.UploadThread)
// 	threadG, uploadCtx := errgroup.NewOrderedGroupWithContext(ctx, thread,
// 		retry.Attempts(3),
// 		retry.Delay(time.Second),
// 		retry.DelayType(retry.BackOffDelay))

// 	for partIndex := range uploadNums {
// 		if utils.IsCanceled(uploadCtx) {
// 			break
// 		}
// 		partIndex := partIndex
// 		partNumber := partIndex + 1 // 分片号从1开始
// 		offset := partIndex * chunkSize
// 		size := min(chunkSize, size-offset)
// 		var reader *stream.SectionReader
// 		var rateLimitedRd io.Reader
// 		sliceMD5 := ""
// 		// 表单
// 		b := bytes.NewBuffer(make([]byte, 0, 2048))
// 		threadG.GoWithLifecycle(errgroup.Lifecycle{
// 			Before: func(ctx context.Context) error {
// 				if reader == nil {
// 					var err error
// 					// 每个分片一个reader
// 					reader, err = ss.GetSectionReader(offset, size)
// 					if err != nil {
// 						return err
// 					}
// 					// 计算当前分片的MD5
// 					sliceMD5, err = utils.HashReader(utils.MD5, reader)
// 					if err != nil {
// 						return err
// 					}
// 				}
// 				return nil
// 			},
// 			Do: func(ctx context.Context) error {
// 				// 重置分片reader位置，因为HashReader、上一次失败已经读取到分片EOF
// 				reader.Seek(0, io.SeekStart)

// 				b.Reset()
// 				w := multipart.NewWriter(b)
// 				// 添加表单字段
// 				err = w.WriteField("preuploadID", createResp.Data.PreuploadID)
// 				if err != nil {
// 					return err
// 				}
// 				err = w.WriteField("sliceNo", strconv.FormatInt(partNumber, 10))
// 				if err != nil {
// 					return err
// 				}
// 				err = w.WriteField("sliceMD5", sliceMD5)
// 				if err != nil {
// 					return err
// 				}
// 				// 写入文件内容
// 				_, err = w.CreateFormFile("slice", fmt.Sprintf("%s.part%d", file.GetName(), partNumber))
// 				if err != nil {
// 					return err
// 				}
// 				headSize := b.Len()
// 				err = w.Close()
// 				if err != nil {
// 					return err
// 				}
// 				head := bytes.NewReader(b.Bytes()[:headSize])
// 				tail := bytes.NewReader(b.Bytes()[headSize:])
// 				rateLimitedRd = driver.NewLimitedUploadStream(ctx, io.MultiReader(head, reader, tail))
// 				// 创建请求并设置header
// 				req, err := http.NewRequestWithContext(ctx, http.MethodPost, uploadDomain+"/upload/v2/file/slice", rateLimitedRd)
// 				if err != nil {
// 					return err
// 				}

// 				// 设置请求头
// 				req.Header.Add("Authorization", "Bearer "+d.AccessToken)
// 				req.Header.Add("Content-Type", w.FormDataContentType())
// 				req.Header.Add("Platform", "open_platform")

// 				res, err := base.HttpClient.Do(req)
// 				if err != nil {
// 					return err
// 				}
// 				defer res.Body.Close()
// 				if res.StatusCode != 200 {
// 					return fmt.Errorf("slice %d upload failed, status code: %d", partNumber, res.StatusCode)
// 				}
// 				var resp BaseResp
// 				respBody, err := io.ReadAll(res.Body)
// 				if err != nil {
// 					return err
// 				}
// 				err = json.Unmarshal(respBody, &resp)
// 				if err != nil {
// 					return err
// 				}
// 				if resp.Code != 0 {
// 					return fmt.Errorf("slice %d upload failed: %s", partNumber, resp.Message)
// 				}

// 				progress := 10.0 + 85.0*float64(threadG.Success())/float64(uploadNums)
// 				up(progress)
// 				return nil
// 			},
// 			After: func(err error) {
// 				ss.FreeSectionReader(reader)
// 			},
// 		})
// 	}

// 	if err := threadG.Wait(); err != nil {
// 		return err
// 	}

// 	return nil
// }

// 上传完毕
// func (d *Open123) complete(preuploadID string) (*UploadCompleteResp, error) {
// 	var resp UploadCompleteResp
// 	_, err := d.Request(UploadComplete, http.MethodPost, func(req *resty.Request) {
// 		req.SetBody(base.Json{
// 			"preuploadID": preuploadID,
// 		})
// 	}, &resp)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &resp, nil
// }

// func (d *Open123) slice(req *UploadDoneReq) error {
// 	r := &BaseResp{
// 		Code: -1,
// 	}
// 	_, err := d.Request(d.qpsInstance[uploadDoneAPI], http.MethodPost, func(req *resty.Request) {
// 		req.SetBody(requ)
// 	}, r)
// 	return err
// }
// func (d *Open123) create(parentFileID int64, filename string, etag string, size int64, duplicate int, containDir bool) (*UploadCreateResp, error) {
// 	var resp UploadCreateResp
// 	_, err := d.Request(UploadCreate, http.MethodPost, func(req *resty.Request) {
// 		req.SetBody(base.Json{
// 			"parentFileId": parentFileID,
// 			"filename":     filename,
// 			"etag":         strings.ToLower(etag),
// 			"size":         size,
// 			"duplicate":    duplicate,
// 			"containDir":   containDir,
// 		})
// 	}, &resp)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &resp, nil
// }

// func (d *Open123) url(preuploadID string, sliceNo int64) (string, error) {
// 	// get upload url
// 	var resp UploadUrlResp
// 	_, err := d.Request(UploadUrl, http.MethodPost, func(req *resty.Request) {
// 		req.SetBody(base.Json{
// 			"preuploadId": preuploadID,
// 			"sliceNo":     sliceNo,
// 		})
// 	}, &resp)
// 	if err != nil {
// 		return "", err
// 	}
// 	return resp.Data.PresignedURL, nil
// }

// func (d *Open123) complete(preuploadID string) (*UploadCompleteResp, error) {
// 	var resp UploadCompleteResp
// 	_, err := d.Request(UploadComplete, http.MethodPost, func(req *resty.Request) {
// 		req.SetBody(base.Json{
// 			"preuploadID": preuploadID,
// 		})
// 	}, &resp)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &resp, nil
// }

// func (d *Open123) async(preuploadID string) (*UploadAsyncResp, error) {
// 	var resp UploadAsyncResp
// 	_, err := d.Request(UploadAsync, http.MethodPost, func(req *resty.Request) {
// 		req.SetBody(base.Json{
// 			"preuploadID": preuploadID,
// 		})
// 	}, &resp)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &resp, nil
// }

// func (d *Open123) Upload(ctx context.Context, file model.FileStreamer, createResp *UploadCreateResp, up driver.UpdateProgress) error {
// 	size := file.GetSize()
// 	chunkSize := createResp.Data.SliceSize
// 	uploadNums := (size + chunkSize - 1) / chunkSize
// 	threadG, uploadCtx := errgroup.NewGroupWithContext(ctx, d.UploadThread,
// 		retry.Attempts(3),
// 		retry.Delay(time.Second),
// 		retry.DelayType(retry.BackOffDelay))

// 	for partIndex := int64(0); partIndex < uploadNums; partIndex++ {
// 		if utils.IsCanceled(uploadCtx) {
// 			return ctx.Err()
// 		}
// 		partIndex := partIndex
// 		partNumber := partIndex + 1 // 分片号从1开始
// 		offset := partIndex * chunkSize
// 		size := min(chunkSize, size-offset)
// 		limitedReader, err := file.RangeRead(http_range.Range{
// 			Start:  offset,
// 			Length: size})
// 		if err != nil {
// 			return err
// 		}
// 		limitedReader = driver.NewLimitedUploadStream(ctx, limitedReader)

// 		threadG.Go(func(ctx context.Context) error {
// 			uploadPartUrl, err := d.url(createResp.Data.PreuploadID, partNumber)
// 			if err != nil {
// 				return err
// 			}

// 			req, err := http.NewRequestWithContext(ctx, "PUT", uploadPartUrl, limitedReader)
// 			if err != nil {
// 				return err
// 			}
// 			req = req.WithContext(ctx)
// 			req.ContentLength = size

// 			res, err := base.HttpClient.Do(req)
// 			if err != nil {
// 				return err
// 			}
// 			_ = res.Body.Close()

// 			progress := 10.0 + 85.0*float64(threadG.Success())/float64(uploadNums)
// 			up(progress)
// 			return nil
// 		})
// 	}

// 	if err := threadG.Wait(); err != nil {
// 		return err
// 	}

// 	uploadCompleteResp, err := d.complete(createResp.Data.PreuploadID)
// 	if err != nil {
// 		return err
// 	}
// 	if uploadCompleteResp.Data.Async == false || uploadCompleteResp.Data.Completed {
// 		return nil
// 	}

// 	for {
// 		uploadAsyncResp, err := d.async(createResp.Data.PreuploadID)
// 		if err != nil {
// 			return err
// 		}
// 		if uploadAsyncResp.Data.Completed {
// 			break
// 		}
// 	}
// 	up(100)
// 	return nil
// }
