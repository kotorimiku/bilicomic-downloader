package main

import (
	"fmt"
	"testing"

	bd "bilicomic-downloader/backend"
)

func TestXxx(t *testing.T) {
	bd.ClientInit("night=0; cf_chl_rc_m=36; jieqiRecentRead=78.4616.0.1.1740640932.0; cf_clearance=VP5cUD99obtGMujKzDNsBsmlZJ3t34XiVqu2qmxL5WQ-1740640933-1.2.1.1-C_aX0gLAa3EFOrv6GFkWkE4RK1C2uvRuI1PGU75TztIPd6F50BEOt36aAFzZZRgxgHI_eMwtZbzfx4is9stVkM4t8YSK0z69QAOB3QkTm.yqyAK1YzwoQIkVaDm996Gwq72Zb.UcHzHYEH.pSwrrSNaMoUN_RpSvoOxmL_I5dkNKeGg_O2jHcBQ83j_jXz3wyOOhNwTVTFu1DoWEblr6mW57JFbrwagijL4bqnrYjtjNBlmV.eKp_7bPd1cLoOT.dC_Tpiqxq.RS2fd3lT5Ah7ht2yu94nI.Ps4LinS10K8gUng6wVVFeYD_7QGZ80r.DWTYBvedM9Q5G_te4fq_5Q")
	// downloader := bd.NewDownloader("159", bd.ConfigInstance)
	// downloader.GetMetadata()
	fmt.Println(bd.GetText("https://www.bilicomic.net/read/78/4616.html"))
}
