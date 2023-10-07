package datawatch

import (
	"context"
	"datamanage/log"
	"github.com/go-mysql-org/go-mysql/replication"
	"os"
	"os/signal"
	"syscall"
)

// watchBinlog 监听binlog日志
func (sdw *SourceDataWatcher) watchBinlog() {
	cfg := replication.BinlogSyncerConfig{
		ServerID:        sdw.ServerID,
		Host:            sdw.Host,
		Port:            sdw.Port,
		User:            sdw.User,
		Password:        sdw.Password,
		Charset:         sdw.Charset,
		RawModeEnabled:  false,
		SemiSyncEnabled: false,
		UseDecimal:      true,
		Logger:          log.GetLogger(),
	}
	syncer := replication.NewBinlogSyncer(cfg)
	position := sdw.getPosition()
	streamer, err := syncer.StartSync(position)
	if err != nil {
		panic(err)
	}
	defer syncer.Close()
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-stopSignal:
			log.Info("Stopping binlog event listener...")
			return
		default:
			ev, err := streamer.GetEvent(context.Background())
			if err != nil {
				log.ErrorF("Failed to get event: %v\n", err)
				continue
			}
			// 分发事件到处理器
			eventType := ev.Header.EventType
			switch e := ev.Event.(type) {
			case *replication.RotateEvent:
				parseEventError(sdw.OnRotate(e))
			case *replication.TableMapEvent:
				parseEventError(sdw.OnTableChanged(e))
			case *replication.QueryEvent:
				parseEventError(sdw.OnDDL(e))
			case *replication.RowsEvent:
				parseEventError(sdw.OnRow(e, eventType))
			}
		}
	}
}

func parseEventError(err error) {
	if err != nil {
		log.Error(err)
	}
}
