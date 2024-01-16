package edgebound

import (
	"net"
	"strconv"
	"time"

	"github.com/singchia/frontier/pkg/repo/dao"
	"github.com/singchia/frontier/pkg/repo/model"
	"github.com/singchia/geminio"
	"k8s.io/klog/v2"
)

func (em *edgeManager) online(end geminio.End) error {
	// TODO transaction
	// cache
	old, ok := em.edges.Swap(end.ClientID(), end)
	if ok {
		// if the old connection exits, offline it
		oldend := old.(geminio.End)
		// we wait the cache and db to clear old end's data
		syncKey := strconv.FormatUint(oldend.ClientID(), 10) + "-" + oldend.RemoteAddr().String()
		sync := em.shub.Add(syncKey)
		if err := oldend.Close(); err != nil {
			klog.Warningf("edge online, kick off old end err: %s, edgeID: %d", err, end.ClientID())
		}
		<-sync.C()
	}

	// memdb
	edge := &model.Edge{
		EdgeID:     end.ClientID(),
		Meta:       string(end.Meta()),
		Addr:       end.RemoteAddr().String(),
		CreateTime: time.Now().Unix(),
	}
	if err := em.dao.CreateEdge(edge); err != nil {
		klog.Errorf("edge online, dao create err: %s, edgeID: %d", err, end.ClientID())
		return err
	}
	return nil
}

func (em *edgeManager) offline(edgeID uint64, addr net.Addr) error {
	// TODO transaction
	legacy := false
	// cache
	value, ok := em.edges.Load(edgeID)
	if ok {
		end := value.(geminio.End)
		if end != nil && end.RemoteAddr().String() == addr.String() {
			legacy = true
			// delete only when the end is old one
			// and the operation should be atomic
			em.edges.CompareAndDelete(edgeID, end)
			klog.Infof("edge offline, edgeID: %d, remote addr: %s", edgeID, end.RemoteAddr().String())
		}
	} else {
		klog.Warningf("edge offline, edgeID: %d not found in cache", edgeID)
	}

	// memdb
	if err := em.dao.DeleteEdge(&dao.EdgeDelete{
		EdgeID: edgeID,
		Addr:   addr.String(),
	}); err != nil {
		klog.Errorf("edge offline, dao delete edge err: %s, edgeID: %d", err, edgeID)
		return err
	}
	if err := em.dao.DeleteEdgeRPCs(edgeID); err != nil {
		klog.Errorf("edge offline, dao delete edge rpcs err: %s, edgeID: %d", err, edgeID)
	}
	if legacy {
		syncKey := strconv.FormatUint(edgeID, 10) + "-" + addr.String()
		em.shub.Done(syncKey)
	}
	return nil
}

// delegations for all ends, called by geminio
func (em *edgeManager) ConnOnline(edgeID uint64, meta []byte, addr net.Addr) error {
	klog.V(4).Infof("edge online, edgeID: %d, meta: %s, addr: %s", edgeID, string(meta), addr.String())
	// notification for others
	return nil
}

func (em *edgeManager) ConnOffline(edgeID uint64, meta []byte, addr net.Addr) error {
	klog.V(4).Infof("edge offline, edgeID: %d, meta: %s, addr: %s", edgeID, string(meta), addr.String())
	return nil
}

func (em *edgeManager) Heartbeat(edgeID uint64, meta []byte, addr net.Addr) error {
	klog.V(5).Infof("edge heartbeat, edgeID: %d, meta: %s, addr: %s")
	return nil
}

func (em *edgeManager) RemoteRegistration(rpc string, edgeID, streamID uint64) {
	klog.V(5).Infof("edge remote rpc registration, rpc: %s, edgeID: %d, streamID: %d", rpc, edgeID, streamID)

	// memdb
	er := &model.EdgeRPC{
		RPC:        rpc,
		EdgeID:     edgeID,
		CreateTime: time.Now().Unix(),
	}
	err := em.dao.CreateEdgeRPC(er)
	if err != nil {
		klog.Errorf("edge remote registration, create edge rpc err: %s, rpc: %s, edgeID: %d, streamID: %d", err, rpc, edgeID, streamID)
	}
}

func (em *edgeManager) GetClientID(meta []byte) (uint64, error) {
	return 0, nil
}