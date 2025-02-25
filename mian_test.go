package main

import (
	"testing"

	bd "bilicomic-downloader/backend"
)

func TestXxx(t *testing.T) {
	bd.ClientInit("night=0; jieqiRecentRead=60.3373.0.1.1740412640.0-159.13203.0.1.1740416624.0; cf_clearance=6Kz4rMl_t4rbmXSHDAYzhHKxLYobH5HqMMOm1JeuhbI-1740416624-1.2.1.1-tRzR8okVuw2PcYZLyMgRFGruVXkUWl75MXfn1LpWHSiN.aqauXu7HHHqTTVtxEEP9MqIxdhvSO81imXizmR.dteRuXwsZwpBgsJfzxOCg8Oq3xEGBtc0MUQUjqLeWT4BQK2L1U0aCRDnuy9psaz3NnFzvy2_VqDiFUFOmzdPIEC5PXDugPRmhHD86bIsGc49rzxWIL2Vb7tWqNwHb98.6daM8LHb8B5Wj5e2jsUrdX.Iuo_lj9QFgiwQQm80AOF4vSx8cTujxPGu8BAuomsHUTMNBWiBzTiuBbrKJyjRVAK_mYYLHXoMdmbUWlS_h1HTPGwUa.dX9JMMjkY0pH8MqA")
	downloader := bd.NewDownloader("159", bd.ConfigInstance)
	downloader.GetMetadata()
	downloader.GetVolume()
	downloader.DownloadList([]int{0})
}
