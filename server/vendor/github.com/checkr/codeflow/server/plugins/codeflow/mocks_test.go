package codeflow

import "encoding/json"

func OktaIdToken() interface{} {
	s := `{"idToken":"eyJhbGciOiJSUzI1NiIsImtpZCI6ImJkTkpwaXU4Uk1SaUJKWUpNZWx1RWpUcnB0ZXZvUTdNOG5tSGdpLVh1MkUifQ.eyJzdWIiOiIwMHVoM2Izd2w4akgzcERTUDF0NSIsIm5hbWUiOiJTYXNvIE1hdGVqaW5hIiwibG9jYWxlIjoiZW4tVVMiLCJlbWFpbCI6InNhc29AY2hlY2tyLmNvbSIsInZlciI6MSwiaXNzIjoiaHR0cHM6Ly9jaGVja3Iub2t0YS5jb20iLCJhdWQiOiJUSnh4MVg2MVJUQ0Y4dXhOcHhsbCIsImlhdCI6MTQ4MjU1MTUyMiwiZXhwIjoxNDgyNTU1MTIyLCJqdGkiOiJJRC5ZYndqeEdLQkZNbml6aS1rM2FvRUxxMlJpLWJaaEk3NlRBdnR3cVdxM1NNIiwiYW1yIjpbInN3ayIsIm1mYSIsInB3ZCJdLCJpZHAiOiIwMG9oM2IzN1FheWtjdm5tcDF0NSIsIm5vbmNlIjoiVTNXdndyQ0RndXVKanNvc1h5YVRFaUNnQVpFRk44QmNmempYcWxjc3RNVldTYXpYZTRvNm84R2kwWGRWWUUzYSIsInByZWZlcnJlZF91c2VybmFtZSI6InNhc29AY2hlY2tyLmNvbSIsImdpdmVuX25hbWUiOiJTYXNvIiwiZmFtaWx5X25hbWUiOiJNYXRlamluYSIsInpvbmVpbmZvIjoiQW1lcmljYS9Mb3NfQW5nZWxlcyIsInVwZGF0ZWRfYXQiOjE0NzkyNTI0NTgsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJhdXRoX3RpbWUiOjE0ODI1NTE1MTMsImdyb3VwcyI6WyJjb2RlZmxvd19kZXZlbG9wZXIiLCJjb2RlZmxvd19hZG1pbiJdfQ.Bb_17q4numVvCGLkJSBuTblM7dXkIde_FDg-PVtr3-RXeIZyLAVHIQjE_jKwwgXuZN78eQ4wUxMCuOBrRvcQk3plO4w-r32m4X03DEgjJZdKALTLvet9PGoQwAF5Qp-3xgKsecTOheSNZEEKP66uvufIbjKlzSY9rkkoHkI0asEnwwyh9xz2xaIi8dTZR1s1-QZkj7K2oGYmDVQRwLfKi6SCF36SY-MfTwfo6FX0LVPEi1pkCkGHCCdA_T2r4Ot8WmSCgl2jQklyyRfN7ah3XmtDiXfaON5sqcgtSDmoGfMTauEJ-vWz8Ates_ghEd-4hbOBclPOJwJFKpbWsO6drw"}`

	var payload interface{}
	_ = json.Unmarshal([]byte(s), &payload)
	return payload
}

func OktaKeys() interface{} {
	s := `{
		"keys": [
		{
			"alg": "RS256",
			"e": "AQAB",
			"n": "kKBxsowLCZyn3bGRi16yGVSF_XYH3EDUVTmagoJkc6ScHBA0d2qTYarFLiPtoJtDD4FNA_IC1PF1RQQuB20m2fm1RKiGb6xckyfxMY1vLzVOvFyt_f41pW-oavrSIFbbc86KNfNvafTOjNw9VLvi29TVVcylDOtDMUsr_6y5PJMG-6q-Xg8BeKlzSqQtVrTuiNvuX9OPVHch4mo51zYMvYdqUniLTijl8AVlW9Tv3wyEXgE6Do9CAudhKRvPUrKKROg4RUmdGuG8-pQVpbL9WrvK51WiMjwnbLOaifMHP6ArkuqD17x0q1PXpSp0tXPg9z905bdeX5Kq5RFS0MNvsQ",
			"kid": "bdNJpiu8RMRiBJYJMeluEjTrptevoQ7M8nmHgi-Xu2E",
			"kty": "RSA",
			"use": "sig"
		},
		{
			"alg": "RS256",
			"e": "AQAB",
			"n": "rXNkvrqqwDdeFWZyZ4PvZcQsRS3yp7bYy3L1ULFQyibk82ZvQZMrK0S20XJDTB9RYJCQruLFLlfCdKtYQobvZw9Ck_m4Wyw1sTgCHCqXAyHmfM62fO9j5WOQNE_Cv2wZI_Fy4ndxNgudJTh2ddrnUXlddDEnzm4wFQV0Qd2Xih6Id2P596pjbvS7Dd93_ELUN4Dwmp-KTKnmkHwLutS6qbgG7iX0ByypN-frb38fwZ-88N_D-7Gwo3T2r__v4PlZRCtjqynft5wwixua6OW-1fRv-A3lUQPGMhAgxYAotdGQvpq8uwUkKCS5FPvtWXNV0-t1RtA_Ft1KOBuaTu8OZw",
			"kid": "S3wRjRulSqL96BYAg9MPa-Obl26aAjuJY-MBRv-mlg4",
			"kty": "RSA",
			"use": "sig"
		}
		]
	}`

	var payload interface{}
	_ = json.Unmarshal([]byte(s), &payload)
	return payload
}
