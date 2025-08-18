package queuehandler

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/OpenListTeam/OpenList/v4/internal/db"
	"github.com/OpenListTeam/OpenList/v4/internal/fs"
	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/OpenListTeam/OpenList/v4/internal/model/tables"
	"github.com/OpenListTeam/OpenList/v4/internal/stream"
	"github.com/OpenListTeam/OpenList/v4/pkg/utils"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// QueueHandlerProxyUpload 代理上传
func QueueHandlerProxyUpload(ctx context.Context, t *asynq.Task) error {
	var err error
	var p tables.SliceUpload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		logrus.Error(t.Payload(), errors.WithStack(err))
		return err
	}
	qiID := t.ResultWriter().TaskID()
	qItem := &tables.QueueItem{
		TaskID: qiID,
	}
	if err := db.GetDb().Where(qItem).First(qItem).Error; err != nil {
		logrus.Error(qiID, errors.WithStack(err))
		return err
	}

	qItem.Status = tables.QueueItemStatusRunning

	defer func() {
		if err != nil {
			qItem.Status = tables.QueueItemStatusFailed
			qItem.Message = err.Error()
		} else {
			qItem.Status = tables.QueueItemStatusSuccess
			qItem.Message = ""
		}
		qItem.EndAt = time.Now()

		//更新
		if err := db.GetDb().Save(&qItem).Error; err != nil {
			logrus.Error(p.ID, errors.WithStack(err))
		}
	}()
	info, err := os.Stat(p.TmpFile)
	if err != nil {
		logrus.Error(p.ID, errors.WithStack(err))
		return err
	}
	f, err := os.Open(p.TmpFile)
	if err != nil {
		logrus.Error(p, errors.WithStack(err))
		return err
	}
	defer f.Close()
	var hashInfo utils.HashInfo
	if p.HashMd5 != "" {
		hashInfo = utils.NewHashInfo(utils.MD5, p.HashMd5)
	}
	if p.HashSha1 != "" {
		hashInfo = utils.NewHashInfo(utils.SHA1, p.HashSha1)
	}

	s := &stream.FileStream{
		Obj: &model.Object{
			Name:     info.Name(),
			Size:     info.Size(),
			Modified: info.ModTime(),
			HashInfo: hashInfo,
		},
		Reader:       f,
		Mimetype:     "application/octet-stream",
		WebPutAsTask: false,
	}

	err = fs.PutDirectly(ctx, p.DstPath, s)
	if err != nil {
		logrus.Error(p, errors.WithStack(err))
		return err
	}
	//移除临时文件
	if err := os.Remove(p.TmpFile); err != nil {
		logrus.Error(p.TmpFile, errors.WithStack(err))
	}

	return nil
}

// QueueHandlerManualMoveCopy 手动移动/复制任务
// func QueueHandlerManualMoveCopy(ctx context.Context, t *asynq.Task) error {
// 	var err error
// 	var p model.Asynqueue
// 	if err := json.Unmarshal(t.Payload(), &p); err != nil {
// 		logrus.Error(t.Payload(), errors.WithStack(err))
// 		return err
// 	}
// 	p.StartTime = time.Now()
// 	p.Status = model.QueueStatusRunning
// 	defer func() {
// 		if err != nil {
// 			p.Status = model.QueueStatusFailed
// 			p.Message = err.Error()
// 		} else {
// 			p.Status = model.QueueStatusSuccess
// 		}
// 		p.EndTime = time.Now()

// 		//更新
// 		if err := db.GetDb().Save(&p).Error; err != nil {
// 			logrus.Error(p.ID, errors.WithStack(err))
// 		}
// 	}()

// 	//移动
// 	if t.Type() == queue.TaskTypeMove {
// 		err = fs.Move(ctx, p.Src, p.Dst)
// 		return err
// 	}
// 	//复制
// 	if t.Type() == queue.TaskTypeCopy {

// 	}

// 	return nil
// }

// // QueueHandlerCronMoveCopy 定时任务移动复制
// func QueueHandlerCronMoveCopy(ctx context.Context, t *asynq.Task) error {
// 	var err error
// 	var p model.Asynqueue
// 	if err := json.Unmarshal(t.Payload(), &p); err != nil {
// 		logrus.Error(t.Payload(), errors.WithStack(err))
// 		return err
// 	}
// 	p.StartTime = time.Now()
// 	p.Status = model.QueueStatusRunning
// 	defer func() {
// 		if err != nil {
// 			p.Status = model.QueueStatusFailed
// 			p.Message = err.Error()
// 		} else {
// 			p.Status = model.QueueStatusSuccess
// 		}
// 		p.EndTime = time.Now()

// 		//更新
// 		if err := db.GetDb().Save(&p).Error; err != nil {
// 			logrus.Error(p.ID, errors.WithStack(err))
// 		}
// 	}()

// 	//移动
// 	if t.Type() == queue.TaskTypeMove {
// 		err = fs.Move(ctx, p.Src, p.Dst)
// 		return err
// 	}
// 	//复制
// 	if t.Type() == queue.TaskTypeCopy {

// 	}

// 	return nil
// }
