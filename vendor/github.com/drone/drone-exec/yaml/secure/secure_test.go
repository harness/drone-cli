package secure

import (
	"testing"

	"github.com/franela/goblin"
)

func Test_Secure(t *testing.T) {

	g := goblin.Goblin(t)
	g.Describe("Secure yaml", func() {

		priv, _ := decodePrivateKey(fakePriv)
		pub := priv.PublicKey
		pem := encodePrivateKey(priv)

		g.It("Should encrypt and decrypt", func() {
			plain := checksumYaml
			encrypted, err := encrypt(plain, &pub)
			g.Assert(err == nil).IsTrue()
			decrypted, err := decrypt(encrypted, priv)
			g.Assert(err == nil).IsTrue()
			g.Assert(plain).Equal(string(decrypted))
		})

		g.It("Should decrypt a yaml with slice parameters", func() {
			secure, err := Parse(sliceEnc, pem)
			g.Assert(err == nil).IsTrue()
			g.Assert(secure.Checksum).Equal("fa4d4048a6bd1a94f2775039ecf29b812d9cfe6b")
			g.Assert(secure.Environment.Map()["FOO"]).Equal("Bar")
			g.Assert(secure.Environment.Map()["BAZ"]).Equal("BOO")
		})

		g.It("Should decrypt a yaml with map parameters", func() {
			secure, err := Parse(mapEnc, pem)
			g.Assert(err == nil).IsTrue()
			g.Assert(secure.Checksum).Equal("fa4d4048a6bd1a94f2775039ecf29b812d9cfe6b")
			g.Assert(secure.Environment.Map()["FOO"]).Equal("Bar")
			g.Assert(secure.Environment.Map()["BAZ"]).Equal("BOO")
		})

		g.It("Should handle an empty environment map", func() {
			secure, err := Parse(checksumEnc, pem)
			g.Assert(err == nil).IsTrue()
			g.Assert(secure.Checksum).Equal("fa4d4048a6bd1a94f2775039ecf29b812d9cfe6b")
			g.Assert(len(secure.Environment.Map())).Equal(0)
		})
	})
}

var sliceYaml = `
checksum: fa4d4048a6bd1a94f2775039ecf29b812d9cfe6b
environment:
  - FOO=BAR
  - BAZ=BOO
`

var mapYaml = `
checksum: fa4d4048a6bd1a94f2775039ecf29b812d9cfe6b
environment:
  FOO: BAR
  BAZ: BOO
`

var checksumYaml = `
checksum: fa4d4048a6bd1a94f2775039ecf29b812d9cfe6b
`

var sliceEnc = `eyJhbGciOiJSU0EtT0FFUCIsImVuYyI6IkExMjhHQ00ifQ.0s31RxVuZUEqEFs79R-1iAR3nN-tt3Zp4ATvsDbBauI5SdQFGzSY_EN58jNghd1CC4RU0Y4ltwoxcwcf8e-buAwnPFlE3xHPlFilNOR-Cs9lrjmaRGbxBOq2P4RRDSlQ6ysTHNyNqHKeZHehQc57om2GAP7Ejxpdqs4OVv_ZhR07wuXw1Meisbh9Mr0DSzGot4Rcv72IqXdEGft1r6SrRasTkJNh1_fYnOSDuiCpUadLEKkA2yePVoln4MxVC4lJPTOcatDEyYvLSHLHc316KNw2hML2hkvwEy90AW_uTbhtmap7q2myRQkhvXdDVuYCEflq2qPTFf_kacgArWnq-g.kh-vFwI-mpv7LzLz.w1i3Db7p7fQAU7u7cPzbUIy0-F1A2lM9afUCC6alyna4yNv2nAcPPBEwxPVnHrn1d-LenWRpcyUgzgTfmkYCODzP_YnjVmgkau2b3sQlK2isoImDAGy5mPo.PuvJO_NeBO07oyylSEX2lw`

var mapEnc = `eyJhbGciOiJSU0EtT0FFUCIsImVuYyI6IkExMjhHQ00ifQ.HzrLexECJiBaP4ptn-uPueg4Irf4H73VcCF4UbiCs0jZur2ss8FJM3OF7ZD-eaSIMPlogpFb5PCQWfr4LMD9fae0ss7ced6Kue1q1wI_sgwUVtDC_kq11yG6qmst4yJ7WCqVnWIEitCov39qyaMsjOLtpTJIuZH3h4Jc19kIIW8Wk_v4Oqd4Yg1JklHRMfrjB3U0OA3U-MLStrzdfk9KzPk_3FQnNXILPAYQZGxh02doLP7sYlo_d8wdV6Km26br4xoW8DK1XUKLvbVE_oXnRcL95eZsvF37ADiHhKI1FfqJo31rxVTSh0HUwdNHOUldtRzyN0oAw4grt_aHGNxvdw.eaJt1GAGBl3voC64.D3IIxEGgkAMo4-UjdYtGax6WFefU5vGqOduDzvU7yOyC_flQ_5C-LrjqC1tPEK9sjdPLt8ghxjtZhCTMFjl1xX77HwAAUNBjyBZVSDaJblsOazkmPZOF.weKJmve-z7DBWMdKLASA3Q`

var checksumEnc = `eyJhbGciOiJSU0EtT0FFUCIsImVuYyI6IkExMjhHQ00ifQ.lduGCINc5DVUD6hi3UHzWNkuKLlsrXudLfrktD3gOI36J6r58DlMAGcUNLfgSAU4v0kf9L407EkB4dtqwXbkhNeihNw69BbYa94QQ1H2uW3BNQbPq-1JeeZLAU6dXmkRXZur4KGNuWls4tMqd-Z9OyRSCBzogDzMf2JGJ-eLSL63zhBCzKGwQ6yE1N6cZsS2NN0-1BEgrAk-dC76motQvcRTHmiosADrEGUAM6xy-LQSkcC8DImpXajv-AFlFv5F4BtFBg9e7MLrVishwZAFKq-lexWLRqlcf7xqgU5GVt6_3VtuoWtVIyUFP4ZnM0KFScKrG6zsd1h7G5_zSf9AMA.YaL_NLt5Ei_7BEeo.7igBRj8A-EfvsT4VafSBCi_68_lelDwcANbtePmZENuuxEaLSRyfMsawKY0Oyc9DdYCsNA.Np3a7xQRMNHK1z6Nb8oGXg`

var fakePriv = `
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA71FaA+otDak2rXF/4h69Tz+OxS6NOWaOc/n7dinHXnlo3Toy
ZzvwweJGQKIOfPNBMncz+8h6oLOByFvb95Z1UEM0d+KCFCCutOeN9NNMw4fkUtSZ
7sm6T35wQUkDOiO1YAGy27hQfT7iryhPwA8KmgZmt7toNNf+WymPR8DMwAAYeqHA
5DIEWWsg+RLohOJ0itIk9q6Us9WYhng0sZ9+U+C87FospjKRMyAinSvKx0Uan4ap
YGbLjDQHimWtimfT4XWCGTO1cWno378Vm/newUN6WVaeZ2CSHcWgD2fWcjFixX2A
SvcvfuCo7yZPUPWeiYKrc5d1CC3ncocu43LhSQIDAQABAoIBAQDIbYKM+sfmxAwF
8KOg1gvIXjuNCrK+GxU9LmSajtzpU5cuiHoEGaBGUOJzaQXnQbcds9W2ji2dfxk3
my87SShRIyfDK9GzV7fZzIAIRhrpO1tOv713zj0aLJOJKcPpIlTZ5jJMcC4A5vTk
q0c3W6GOY8QNJohckXT2FnVoK6GPPiaZnavkwH33cJk0j1vMsbADdKF7Jdfq9FBF
Lx+Za7wo79MQIr68KEqsqMpmrawIf1T3TqOCNbkPCL2tu5EfoyGIItrH33SBOV/B
HbIfe4nJYZMWXhe3kZ/xCFqiRx6/wlc5pGCwCicgHJJe/l8Y9OticDCCyJDQtD8I
6927/j2NAoGBAPNRRY8r5ES5f8ftEktcLwh2zw08PNkcolTeqsEMbWAQspV/v+Ay
4niEXIN3ix2yTnMgrtxRGO7zdPnMaTN8E88FsSDKQ97lm7m3jo7lZtDMz16UxGmd
AOOuXwUtpngz7OrQ25NXhvFYLTgLoPsv3PbFbF1pwbhZqPTttTdg5so3AoGBAPvK
ta/n7DMZd/HptrkdkxxHaGN19ZjBVIqyeORhIDznEYjv9Z90JvzRxCmUriD4fyJC
/XSTytORa34UgmOk1XFtxWusXhnYqCTIHG/MKCy9D4ifzFzii9y/M+EnQIMb658l
+edLyrGFla+t5NS1XAqDYjfqpUFbMvU1kVoDJ/B/AoGBANBQe3o5PMSuAD19tdT5
Rnc7qMcPFJVZE44P2SdQaW/+u7aM2gyr5AMEZ2RS+7LgDpQ4nhyX/f3OSA75t/PR
PfBXUi/dm8AA2pNlGNM0ihMn1j6GpaY6OiG0DzwSulxdMHBVgjgijrCgKo66Pgfw
EYDgw4cyXR1k/ec8gJK6Dr1/AoGBANvmSY77Kdnm4E4yIxbAsX39DznuBzQFhGQt
Qk+SU6lc1H+Xshg0ROh/+qWl5/17iOzPPLPXb0getJZEKywDBTYu/D/xJa3E/fRB
oDQzRNLtuudDSCPG5wc/JXv53+mhNMKlU/+gvcEUPYpUgIkUavHzlI/pKbJOh86H
ng3Su8rZAn9w/zkoJu+n7sHta/Hp6zPTbvjZ1EijZp0+RygBgiv9UjDZ6D9EGcjR
ZiFwuc8I0g7+GRkgG2NbfqX5Cewb/nbJQpHPO31bqJrcLzU0KurYAwQVx6WGW0He
ERIlTeOMxVo6M0OpI+rH5bOLdLLEVhNtM/4HUFi1Qy6CCMbN2t3H
-----END RSA PRIVATE KEY-----
`
