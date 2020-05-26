package main

type (
	WebSocketMsg struct {
		Type byte        `json:"type"`
		Data interface{} `json:"data"`
	}

	NodeStatus struct {
		Mem *MemMetrics `json:"mem"`
	}

	MemMetrics struct {
		Sys          uint64 `json:"sys"`
		HeapSys      uint64 `json:"heap_sys"`
		HeapInuse    uint64 `json:"heap_inuse"`
		HeapIdle     uint64 `json:"heap_idle"`
		HeapReleased uint64 `json:"heap_released"`
		HeapObjects  uint64 `json:"heap_objects"`
		MSpanInuse   uint64 `json:"m_span_inuse"`
		MCacheInuse  uint64 `json:"m_cache_inuse"`
		StackSys     uint64 `json:"stack_sys"`
		NumGC        uint32 `json:"num_gc"`
		LastPauseGC  uint64 `json:"last_pause_gc"`
	}

	TPSMetrics struct {
		Incoming uint32 `json:"incoming"`
		New      uint32 `json:"new"`
		Outgoing uint32 `json:"outgoing"`
	}
)

const (
	// MsgTypeNodeStatus is the type of the NodeStatus message.
	MsgTypeNodeStatus byte = iota
	// MsgTypeTPSMetric is the type of the transactions per second (TPS) metric message.
	MsgTypeTPSMetric
	// MsgTypeTipSelMetric is the type of the TipSelMetric message.
	MsgTypeTipSelMetric
	// MsgTypeTx is the type of the Tx message.
	MsgTypeTx
	// MsgTypeMs is the type of the Ms message.
	MsgTypeMs
	// MsgTypePeerMetric is the type of the PeerMetric message.
	MsgTypePeerMetric
	// MsgTypeConfirmedMsMetrics is the type of the ConfirmedMsMetrics message.
	MsgTypeConfirmedMsMetrics
	// MsgTypeVertex is the type of the Vertex message for the visualizer.
	MsgTypeVertex
	// MsgTypeSolidInfo is the type of the SolidInfo message for the visualizer.
	MsgTypeSolidInfo
	// MsgTypeConfirmedInfo is the type of the ConfirmedInfo message for the visualizer.
	MsgTypeConfirmedInfo
	// MsgTypeMilestoneInfo is the type of the MilestoneInfo message for the visualizer.
	MsgTypeMilestoneInfo
	// MsgTypeTipInfo is the type of the TipInfo message for the visualizer.
	MsgTypeTipInfo
	// MsgTypeDatabaseSizeMetric is the type of the database Size message for the metrics.
	MsgTypeDatabaseSizeMetric
	// MsgTypeDatabaseCleanupEvent is the type of the database cleanup message for the metrics.
	MsgTypeDatabaseCleanupEvent
)
