package zfs

/*
 * Source: include/sys/fs/zfs.h
 * The following are configuration names used in the nvlist describing a pool's
 * configuration. New on-disk names should be prefixed with "<reversed-DN_s>:"
 * (e.g. "org.openzfs:") to avoid conflicting names being developed
 * independently.
 */
const (
	PoolConfigVersion          = "version"
	PoolConfigPoolName         = "name"
	PoolConfigPoolState        = "state"
	PoolConfigPoolTXG          = "txg"
	PoolConfigPoolGUID         = "pool_guid"
	PoolConfigCreateTXG        = "create_txg"
	PoolConfigTopGUID          = "top_guid"
	PoolConfigVdevTree         = "vdev_tree"
	PoolConfigType             = "type"
	PoolConfigChildren         = "children"
	PoolConfigID               = "id"
	PoolConfigGUID             = "guid"
	PoolConfigIndirectObject   = "com.delphix:indirect_object"
	PoolConfigIndirectBirths   = "com.delphix:indirect_births"
	PoolConfigPrevIndirectVdev = "com.delphix:prev_indirect_vdev"
	PoolConfigPath             = "path"
	PoolConfigDevId            = "devid"
	PoolConfigMetaslabArray    = "metaslab_array"
	PoolConfigMetaslabShift    = "metaslab_shift"
	PoolConfigAShift           = "ashift"
	PoolConfigASize            = "asize"
	PoolConfigDtl              = "dtl"
	PoolConfigScanStats        = "scan_stats"       /* not stored on disk */
	PoolConfigRemovalStats     = "removal_stats"    /* not stored on disk */
	PoolConfigCheckpointStats  = "checkpoint_stats" /* not on disk */
	PoolConfigVdevStats        = "vdev_stats"       /* not stored on disk */
	PoolConfigIndirectSize     = "indirect_size"    /* not stored on disk */

	/* container nvlist of extended stats */
	PoolConfigVdevStatsEx = "vdev_stats_ex"

	/* active queue read/write stats */
	PoolConfigVdevSyncRActiveQueue  = "vdev_sync_r_active_queue"
	PoolConfigVdevSyncWActiveQueue  = "vdev_sync_w_active_queue"
	PoolConfigVdevAsyncRActiveQueue = "vdev_async_r_active_queue"
	PoolConfigVdevAsyncWActiveQueue = "vdev_async_w_active_queue"
	PoolConfigVdevScrubActiveQueue  = "vdev_async_scrub_active_queue"
	PoolConfigVdevTrimActiveQueue   = "vdev_async_trim_active_queue"

	/* queue sizes */
	PoolConfigVdevSyncRPendQueue  = "vdev_sync_r_pend_queue"
	PoolConfigVdevSyncWPendQueue  = "vdev_sync_w_pend_queue"
	PoolConfigVdevAsyncRPendQueue = "vdev_async_r_pend_queue"
	PoolConfigVdevAsyncWPendQueue = "vdev_async_w_pend_queue"
	PoolConfigVdevScrubPendQueue  = "vdev_async_scrub_pend_queue"
	PoolConfigVdevTrimPendQueue   = "vdev_async_trim_pend_queue"

	/* latency read/write histogram stats */
	PoolConfigVdevTotRLatHisto   = "vdev_tot_r_lat_histo"
	PoolConfigVdevTotWLatHisto   = "vdev_tot_w_lat_histo"
	PoolConfigVdevDiskRLatHisto  = "vdev_disk_r_lat_histo"
	PoolConfigVdevDiskWLatHisto  = "vdev_disk_w_lat_histo"
	PoolConfigVdevSyncRLatHisto  = "vdev_sync_r_lat_histo"
	PoolConfigVdevSyncWLatHisto  = "vdev_sync_w_lat_histo"
	PoolConfigVdevAsyncRLatHisto = "vdev_async_r_lat_histo"
	PoolConfigVdevAsyncWLatHisto = "vdev_async_w_lat_histo"
	PoolConfigVdevScrubLatHisto  = "vdev_scrub_histo"
	PoolConfigVdevTrimLatHisto   = "vdev_trim_histo"

	/* request size histograms */
	PoolConfigVdevSyncIndRHisto  = "vdev_sync_ind_r_histo"
	PoolConfigVdevSyncIndWHisto  = "vdev_sync_ind_w_histo"
	PoolConfigVdevAsyncIndRHisto = "vdev_async_ind_r_histo"
	PoolConfigVdevAsyncIndWHisto = "vdev_async_ind_w_histo"
	PoolConfigVdevIndScrubHisto  = "vdev_ind_scrub_histo"
	PoolConfigVdevIndTrimHisto   = "vdev_ind_trim_histo"
	PoolConfigVdevSyncAggRHisto  = "vdev_sync_agg_r_histo"
	PoolConfigVdevSyncAggWHisto  = "vdev_sync_agg_w_histo"
	PoolConfigVdevAsyncAggRHisto = "vdev_async_agg_r_histo"
	PoolConfigVdevAsyncAggWHisto = "vdev_async_agg_w_histo"
	PoolConfigVdevAggScrubHisto  = "vdev_agg_scrub_histo"
	PoolConfigVdevAggTrimHisto   = "vdev_agg_trim_histo"

	/* number of slow ios */
	PoolConfigVdevSlowIos = "vdev_slow_ios"

	PoolConfigVdevEncSysfsPath = "vdev_enc_sysfs_path"

	PoolConfigWholeDisk       = "whole_disk"
	PoolConfigErrcount        = "error_count"
	PoolConfigNotPresent      = "not_present"
	PoolConfigSpares          = "spares"
	PoolConfigIsSpare         = "is_spare"
	PoolConfigNparity         = "nparity"
	PoolConfigHostid          = "hostid"
	PoolConfigHostname        = "hostname"
	PoolConfigLoadedTime      = "initial_load_time"
	PoolConfigUnspare         = "unspare"
	PoolConfigPhysPath        = "phys_path"
	PoolConfigIsLog           = "is_log"
	PoolConfigL2cache         = "l2cache"
	PoolConfigHoleArray       = "hole_array"
	PoolConfigVdevChildren    = "vdev_children"
	PoolConfigIsHole          = "is_hole"
	PoolConfigDDTHistogram    = "ddt_histogram"
	PoolConfigDDTObjStats     = "ddt_object_stats"
	PoolConfigDDTStats        = "ddt_stats"
	PoolConfigSplit           = "splitcfg"
	PoolConfigOrigGUID        = "origGUI_d"
	PoolConfigSplitGUID       = "splitGUI_d"
	PoolConfigSplitList       = "guid_list"
	PoolConfigRemoving        = "removing"
	PoolConfigResilverTXG     = "resilver_txg"
	PoolConfigRebuildTXG      = "rebuild_txg"
	PoolConfigComment         = "comment"
	PoolConfigSuspended       = "suspended"        /* not stored on disk */
	PoolConfigSuspendedReason = "suspended_reason" /* not stored */
	PoolConfigTimestamp       = "timestamp"        /* not stored on disk */
	PoolConfigBootfs          = "bootfs"           /* not stored on disk */
	PoolConfigMissingDevices  = "missing_vdevs"    /* not stored on disk */
	PoolConfigLoadInfo        = "load_info"        /* not stored on disk */
	PoolConfigRewindInfo      = "rewind_info"      /* not stored on disk */
	PoolConfigUnsupFeat       = "unsup_feat"       /* not stored on disk */
	PoolConfigEnabledFeat     = "enabled_feat"     /* not stored on disk */
	PoolConfigCanRdonly       = "can_rdonly"       /* not stored on disk */
	PoolConfigFeaturesForRead = "features_for_read"
	PoolConfigFeatureStats    = "feature_stats" /* not stored on disk */
	PoolConfigErrata          = "errata"        /* not stored on disk */
	PoolConfigVdevTopZap      = "com.delphix:vdev_zap_top"
	PoolConfigVdevLeafZap     = "com.delphix:vdev_zap_leaf"
	PoolConfigHasPerVdevZaps  = "com.delphix:has_per_vdev_zaps"
	PoolConfigResilverDefer   = "com.datto:resilver_defer"
	PoolConfigCachefile       = "cachefile"      /* not stored on disk */
	PoolConfigMmpState        = "mmp_state"      /* not stored on disk */
	PoolConfigMmpTXG          = "mmp_txg"        /* not stored on disk */
	PoolConfigMmpSeq          = "mmp_seq"        /* not stored on disk */
	PoolConfigMmpHostname     = "mmp_hostname"   /* not stored on disk */
	PoolConfigMmpHostID       = "mmp_hostid"     /* not stored on disk */
	PoolConfigAllocationBias  = "alloc_bias"     /* not stored on disk */
	PoolConfigExpansionTime   = "expansion_time" /* not stored */
	PoolConfigRebuildStats    = "org.openzfs:rebuild_stats"

	/*
	 * the persistent vdev state is stored as separate values rather than a single
	 * 'vdevState' entry.  this is because a device can be in multiple states, such
	 * as offline and degraded.
	 */
	PoolConfigOffline  = "offline"
	PoolConfigFaulted  = "faulted"
	PoolConfigDegraded = "degraded"
	PoolConfigRemoved  = "removed"
	PoolConfigFru      = "fru"
	PoolConfigAuxState = "aux_state"
)
