package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	arnList = []string{
		"24234261344000020134087",
		"24234261344000020133907",
		"24234261344000020135332",
	}

	orderIDList = []string{
		"CAQACAQNryLheQlKuVoPlBrZUKB9-l_UQ-VLckQK92NKb797g0gBW2Nxe-Qmpd-cq001c001",
		"CAQACAQM9elMUIRwgo2fNtSB9ZUu_7sW90phoGRBu3whCdkuWRACj32gU0s7TMjgq001c001",
		"CAQACAQOrwV1Vd4f48lD_ynf4wJrhPAe2gUweDpsDe7xLLaixBwDyex5VgdWKUFsq001c001",
		"CAQACAQM3rPlOsvbKazadx8Bh0XZGvjSosP1xuhMQe16y5aQrvABre3FOsLe6YGEq001c001",
		"CAQACAQOkM9tW1ElQINm-FmrqvEiUL5eu_sHEWemF7SYXa0RXxAAg7cRW_muH3jIq001c001",
		"CAQACAQOJqL945iN0gnA89ndW8WLz-OHlKXA1e_vpdOUNKpsfywCCdDV4KR4quowq001c001",
		"CAQACAQOorEVEOX5pFcZYHJcZClYerMKUhoWBiO7DPAuBXgFvwwAVPIFEhi0Jr0Mq001c001",
		"CAQACAQOldkLr0RxvPqrhwZuoxUtspRua3WQZ2LboEEHTo-Pq4QA-EBnr3eb6Ndgq001c001",
		"CAQACAQNyLaBK9ItM_AmA-4bIR6jLnJTvhSNCfolApr9qFkLhswD8pkJKhZCxTF4q001c001",
		"CAQACAQMF0mUlI1M9KNKYYt2agV522Im27kPHhaPOL6AICmSqPAAoL8cl7mhr7pgq001c001",
		"CAQACAQOtQxo7UrnT43BFE4Kd_W_K2-Vq4URXdbKXB458zOFJuwDjB1c74bzeBuoq001c001",
		"CAQACAQPfwpwT6gW-2tfB7YkD7pjqA6gVBsxTBj87_rRJtAV0tgDa_lMTBs1GTBoq001c001",
		"CAQACAQO6t0W-8BhPJ2zwfNJc1QtVYU3XhCSD24nFNq2cxNVBDwAnNoO-hHD-Khkq001c001",
		"CAQACAQOhZnTw6853tyNJkNz72PNImeDE9Ioymnz7wQAdu6zZqAC3wTLw9DvmXy8q001c001",
		"CAQACAQPDalRMSQlQCAermsVnX_7xn7eFE1oM1k9opWrhNd6afwAIpQxME3x11JEq001c001",
		"CAQACAQPvSSHfBDbomG79rxkepJbay7IKBYZbngooRTBXhqbrqgCYRVvfBZXjb8Aq001c001",
		"CAQACAQMnBV4dYxQdDFYTCzXM-aCbCFM9Xain6JL_oaGibHpEjAAMoacdXfxz-Foq001c001",
		"CAQACAQMCwCH_E7QcofOf-m50tjKboRuPJyzvJeXIUeeXL2TuYgChUe__J50QiGkq001c001",
		"CAQACAQPxpDXeERX0uksCf7Y4J4naj1roEnKSGEELmF7kAK6BcAC6mJLeEjvvFsIq001c001",
		"CAQACAQMGGSANOm5aYghjQVHkk_zYN08oFm3srC3FlfsOpHg0qwBilewNFuljVvoq001c001",
		"CAQACAQNqpxnloxYC-R1n4WdXdaMn-3G6WIVupqsK71AYiyq55AD5727lWFDdqJUq001c001",
		"CAQACAQPQpHAy116bTAz0PsvBNQ-dYRj8fSd4uqeLgYPIgFokOgBMgXgyfcroTaQq001c001",
		"CAQACAQMonbvzZHVhoMM19tciic3dkRZc1p3ryeVenpjtVwyphgCgnuvz1s-kd14q001c001",
		"CAQACAQNlcMhlHB258PabYWtFoH9rtgsKkcZrJjXnEatZTqHLSgDwEWtlkSFiq0cq001c001",
		"CAQACAQO10B0AzQ-WwZaZDGswy3IFBo6FPZ3ETGpm2BFH2ElkRADB2MQAPej8tqYq001c001",
		"CAQACAQNEvKa7tY-NK_px3hSYg151MXC2HyPZ6nMjGfPy83xF1AArGdm7H3muha0q001c001",
		"CAQACAQNUoOlLoiAQeVtYbot8jQ8K4bAgk47TDNuZqlWooqfr2wB5qtNLk0BknE0q001c001",
		"CAQACAQP8Vk6PfC85DcyzkU55qJhrjQrGcN9QlK3oIHwyBuGC1gANIFCPcEhtAOYq001c001",
		"CAQACAQMuYhdL1-lYEstRqOK5e78LQLKClvZDOssDj7Z0OjcGmAASj0NLls7c3kAq001c001",
		"CAQACAQN3d8fx7-B9qaFFukWsYrDgzcFdXe3L2ND1tuP2Hl2wggCptsvxXSfr1iAq001c001",
		"CAQACAQO9EhDFl9dzMHvajXSRQHO3z4ZoKSpzG9jKQs9xn37nygAwQnPFKYNMhs0q001c001",
		"CAQACAQO8T42YqK4sLJDFgrhELpCUYdzhN1lD2Y5-R4EfgaVNUAAsR0OYN-jfJUoq001c001",
		"CAQACAQOTF-LyDSFH8olD0EoyjGmISIl80Ns09CgZFjh-hv7YZwDyFjTy0H6Baqoq001c001",
		"CAQACAQMQD5n6OdG9FAtb-kabADtWV57QIzxmivPxJtLMR7clYwAUJmb6I3GhWN8q001c001",
		"CAQACAQNDjEohBFqPoncr7kKPsj5Ve9pJzSIRDUBw5rYW0L01iACi5hEhzaYkK28q001c001",
		"CAQACAQPCmc_YKDmc-IDgJIiD3GIcHHcrDOs-vQFCnNj_eOkGfwD4nD7YDPjp4o0q001c001",
		"CAQACAQPpMcrfWOhti6SSfV3ooPdb2TYIhXhLc8YzqkQ8chg3aQCLqkvfhZHY88gq001c001",
		"CAQACAQMNfaLPtoJwVoICcaZJfQ6KQ8z92vPCLYUO-gFkoUXnRQBW-sLP2nBbq3oq001c001",
		"CAQACAQN4ZKTFQ3NTcyI1DQ3Piq-qmdYKOMKk3nI-zTUkW1WV8ABzzaTFOD9xUYgq001c001",
		"CAQACAQNqkLkjjPy0op3izvUqc6QSNLjkrKdrHKV5JkJP8z8T5QCiJmsjrLVnK6Qq001c001",
		"CAQACAQPfxWz4cbjS7VhrLttDQua5ksHZRjAVRtkyzUrODqdDyADtzRX4Rg1wus8q001c001",
		"CAQACAQPiBY38miiSySpAotoybnaQQ3Ta5Gzt9anJ_CO56p9UqwDJ_O385OS_eLQq001c001",
		"CAQACAQP2uWbfjzfsT81UbKsiu-l11FJj3MIkNtgGGw9o6g5HXgBPGyTf3BsJL5Yq001c001",
		"CAQACAQP9wnhDw3PPIReQ4U5E3KrQA0b58UsG4z4OR_JgCXJOdgAhRwZD8fjoOFUq001c001",
		"CAQACAQOIaiv907kzQjeM1c43UOeoyKslVGKxDc3OSjFBKWKBYQBCSrH9VBDq56Aq001c001",
		"CAQACAQN097KBCdwyI5BjXGL5r3ezfKFR38noFwMYglMaBMNACwAjguiB3zRRjXQq001c001",
		"CAQACAQMERbVy_eGIsVjtRhuv22wJgPFc1-GuV2BacqQ69gG8sgCxcq5y1xpOz8Iq001c001",
		"CAQACAQMUwIno-9fKny3tKxeddix8anQHwZKheCzfNRs03fKlUQCfNaHowQ8WssAq001c001",
		"CAQACAQNZObWEHCxvuezfLvYA1xNHzmmGblBbFnEB0wCcA3gLnQC501uEboD-C6Mq001c001",
		"CAQACAQOjpK7z3gopY04cwCWXXVb6tdIQAj72I9yhEWV6c9I5wABjEfbzArZV5-Mq001c001",
		"CAQACAQPMB5-p1omKlMtDg_B2vTDCqmBsuJVCyfl9iOkRTAJYwgCUiEKpuNBzfOEq001c001",
		"CAQACAQOD9rUpFjn-RORa7Ztcd-yMvAA5Cts2ZYKip2d8BE_9cQBEpzYpClYBRbEq001c001",
		"CAQACAQOutmFEqqcW85mqrmm8UxL5TiLM2-cs1CKEhlPIKqJoyADzhixE2-oXLBAq001c001",
		"CAQACAQPr0Ovseqb_l8CY8K2Pv2HA6iPDXJ_a7PVJyTmc186JQgCXydrsXHlMz-Uq001c001",
		"CAQACAQPUMupRcHK5bMM9W9uynvn2kv-GknAsE8rmFotHrTh-YwBsFixRkv80ve0q001c001",
		"CAQACAQNgkuAG64e490IDJpUca9uzwiZQIAGjixC2mTGdswLCqgD3maMGIN0aF7kq001c001",
		"CAQACAQNRh6ryrUUKKAvtx-4uoFyuuqYe0Db-hqdajMjoEmYOwAAojP7y0CPWJtEq001c001",
		"CAQACAQN15dyZok7Rh2tyOn4-dMAPM2AZtvMoS07l8jYSZzfPQwCH8iiZtpzsv9oq001c001",
		"CAQACAQP749NyJCRH9To43NWxA33uRFsiihLPQHpP0YgGQgx8FAD10c9yilKNASQq001c001",
		"CAQACAQMSaAEWqLR9AeTRy-n-Eim0UfL6V74ZRZ3iVrlPbCBnvAABVhkWVxnCpdMq001c001",
		"CAQACAQO-7Gq98MseyM9tc2SPK7XXdp0rQ4PLhWY3cjf3m2wZogDIcsu9Q33vAwwq001c001",
		"CAQACAQPu6Ih9fW_1SozBAvGLhc2HxNB_ceWXP1umcRKB_Zd0rgBKcZd9cZpOC18q001c001",
		"CAQACAQOd1K2fzYdD4IlJbbT80C4ldOorU-7KNJO_c9TWreflTgDgc8qfU1uzezwq001c001",
		"CAQACAQNz3MBGMMrEuO1NQNTlBd6KIJuC8jHjp10px-CJ5bYloQC4x-NG8lpj0JEq001c001",
		"CAQACAQPpv84lxUAimnRYXQiceRw1-Pc2nnknvWguj0Cg3cFYiwCajyclniQ1pR8q001c001",
		"CAQACAQPtSaRYxPL0ROyYEPr2RbHNMmf-CkMeSme4AQuTXFjaSwBEAR5YChWfQKgq001c001",
		"CAQACAQMtjGOdFkdILHnK3ui-m0pQ-9CyMnsLV3gE49MMnMFP5wAs4wudMmXRi04q001c001",
		"CAQACAQPMQtzF1yPD8ftNrmeLRSRjrD1OFnsiGpXpCS-GkPfQbgDxCSLFFlovVqsq001c001",
		"CAQACAQOuwyowP0E7BXaR2aePpUMlllSi-oPURkQWW8-ZQwXPxgAFW9Qw-ruorn4q001c001",
		"CAQACAQO4DN9MhmJfpoJyDBcnKxPszC5Kh7fUcUzxlYePb7By9ACmldRMh5vB-FIq001c001",
		"CAQACAQPE92X70HKbMXtssarlVc-pGuVyKiVo-TnSFU7SYfRivAAxFWj7KkJbtpQq001c001",
		"CAQACAQNF1gfGrqwZ9fj6ZuAz4XTasSm6RRCpO4zhi9AFA-RqKgD1i6nGRUja_ewq001c001",
		"CAQACAQN2ha82j0OyCOcT3eVYUJU3UX4iwYSOWVQDWb0scLSTMgAIWY42wUIrowUq001c001",
		"CAQACAQN4tP9-VdA0jiVL4rAbuOEuR6NTKFje4WfLAXel5aa-5ACOAd5-KHsRWWAq001c001",
		"CAQACAQNKSurLBc2gTanAX-AQwR47niLOuiVeY8Jw1u0H4D4LBABN1l7Lui8xNWgq001c001",
		"CAQACAQPA7SZmyqoF-XSql3WSWQ-K_dNdN9aBTarMXST4BcaRfAD5XYFmN_CGVj4q001c001",
		"CAQACAQP_O4doSEvyPYW28Ue99bgT9OQ3hHnlVvHRoJZ4vve_4wA9oOVohA2U0f4q001c001",
		"CAQACAQMX05mFUDrgOWIxzjOpMoOlP6OQBl3hmwmYcqqjEZwyJgA5cuGFBpJj40cq001c001",
		"CAQACAQMFwcDoWsqRCsiYAWEPIerP6DZa2u-Nuj3GPsczju_I1gAKPo3o2nqz-F4q001c001",
		"CAQACAQMd5krNBI4tkQjWW87ZhAoO5sVUibiO--Snwx1RfVl35QCRw47NiT6EkiAq001c001",
		"CAQACAQMZAmBhyAEa3cAzjmE_-vssS6MYslb30mNTh9j_Z_VvuwDdh_dhsovpEhsq001c001",
		"CAQACAQPwIOpMABD056KXbegT_rAZ-ZUvoJoKsjaG7i5-akrNYQDn7gpMoHF3Gw0q001c001",
		"CAQACAQM3784RnMHtLExsu3tqH0S8wx3iaQEvcjbGvBOjEr2fiQAsvC8RaXkLi7Eq001c001",
		"CAQACAQOO1lLeQrXFWdekhH_lLVjWbfNcSPemc7iun_0KE70KVABZn6beSPqLEyYq001c001",
		"CAQACAQMBl_7Hdxumv4agDLAvBTb6XC_iymlTsaW4iA3BQlZ2eQC_iFPHyl2Sa1kq001c001",
		"CAQACAQPQop1R2fC9qodQwl64xK5kZI4p_AFp0imClIKf-RqmUQCqlGlR_J_RWjQq001c001",
		"CAQACAQNyB2CC5y_au9hkHTq4nuxi0hCeWC_-VDrtQfCzn80-VAC7Qf6CWBtBp3Mq001c001",
		"CAQACAQNcvdDwa9JVvwfexNHpLWKQrm0Xur5B9M0Pyf5BdI1CtAC_yUHwuqCvdsgq001c001",
		"CAQACAQO6S_0q9pexLYx2FsqvLRfRAtkkaiKqQqeBimyM9CgkoAAtiqoqamLpDKcq001c001",
		"CAQACAQPPpNdBgrCA4Meu6HScdBu9Enh9xUoEdJBSZGnUmLkNMgDgZARBxVwfLK8q001c001",
		"CAQACAQNhnjNTV2zuWOS9lFUpaQgaWQxO6M07oNSqoOgZvs6dUwBYoDtT6IUr6-Iq001c001",
		"CAQACAQNi0dC6aSb2M79gg6ECLLY-ExVYVke5Fpxz7dWylxp1EAAz7bm6VlKB_w4q001c001",
		"CAQACAQNssty6QVvp1vZFsD-LkXMyLhH2_nCa54ZogbTxyjDSMgDWgZq6_lswT2Uq001c001",
		"CAQACAQMUvOMmJCNOduCxt3S1ev-Hn-aE8dh17cASPdi-QxcnqAB2PXUm8Sefcb4q001c001",
		"CAQACAQPfVZgaWaEVe-xEbEBx3NiNJ3qGElQIqsvCbIIeSULV8wB7bAgaEt5y4T0q001c001",
		"CAQACAQO5se7CHKXOw6GTPUQow6F5UOexMGj_4tBujG4TaEYJ3QDDjP_CMMLXxdMq001c001",
		"CAQACAQPmkyMB6yhdDAHFXQuLRmr5rhH5rlVD8mwtwx9F_d47TQAMw0MBrlBFW7wq001c001",
		"CAQACAQPI636hdxavM1E0UCZWIrcDiVs4WVFl4UB2Ry5vK9hqsgAzR2WhWbkPU5sq001c001",
		"CAQACAQNXFf3bYiaADTopt3YZ_Mv3wsrNiW8tVe9iXoIB8z8bvgANXi3biZlhh7Qq001c001",
		"CAQACAQPBtpeAGiJ3q7DmiN2LSEulndqXcFfWKvshV-buV3GpSgCrV9aAcEvFr7kq001c001",
		"CAQACAQOnK6jsetQdVIn33URU3MGDIAlvZ3IqTHQqwcL18P1J0wBUwSrsZ5YTC6Eq001c001",
		"CAQACAQNBQGlDEQwQLmbXHjmSnJ5FIXHHRikEBrdQpXKMHlRKogAupQRDRnYEAVgq001c001",
		"CAQACAQMNLJ6r6L7f8LhljxdkBkoLs1uua49HC8yHlbWXTzIh5QDwlUeraxuk9U4q001c001",
		"CAQACAQNiOLzmktk6QMkD0AteaCyGldmYwrv297XUUvfnkgr7rABAUvbmwhK998oq001c001",
		"CAQACAQOF7I11aP0cgDxMZlJ-7uRP26uZ8mn-oLVZrfTK5h0c9ACArf518nO2uOwq001c001",
		"CAQACAQPNxdLTPB0LOdLSxwITE4Y1PzLsxKF1Mtgu1RZ3KmzbrwA51XXTxNcH2xwq001c001",
		"CAQACAQPmxqaASEY8Erkr7snJymBgmGHzMy1qv5JoqUx9EWXL1QASqWqAMx0_yloq001c001",
		"CAQACAQMega1TZMZgXM19TGEzsOD7Z8ePrhTZmm97OF5FXOEAwwBcONlTrto5m8sq001c001",
		"CAQACAQPUqPvlmFyu37LX3_UKSImUjJnXYgENxs-0qtnwSfRAbQDfqg3lYq6sTiEq001c001",
		"CAQACAQOn8SaTU6ZknGkn_LlLlR76adl5jN9UlKjyj3c17e6JigCcj1STjE4UwMAq001c001",
		"CAQACAQOj5GFKEpLj4InkCKOvQtIrqUW09Vu52e3e4QOJDxiY7gDg4blK9ZUvOnMq001c001",
		"CAQACAQOQhjGb0jisqXvCNN7Wsd85dreYTD3XzvYoDky24tZfngCpDtebTKxNo1gq001c001",
		"CAQACAQOYWoWZP3ie9jPMwRZIk04ELsM1qW0EFMDPWHzL7n_kHwD2WASZqbNTHK0q001c001",
		"CAQACAQPaB6LBZ0QIRFPSTu9sY24tmXA36E8GvgV4-Th_lUpHggBE-QbB6LA4s6Aq001c001",
		"CAQACAQPUA5obJGWQ7cGJtuvhsXJRa4dcqc7WW5jbDvykTHiqXgDtDtYbqSy9NNAq001c001",
		"CAQACAQOrzHp_f6DtMoehhZogr_IxS67q-noPvROh37cBvZC18AAy3w9_-ht2Cywq001c001",
		"CAQACAQNmg3JvKHQstxCp63680NkQLwnuJwxOfijyPos0mRiRWQC3Pk5vJ_933vkq001c001",
		"CAQACAQON2BWy_riipL6Dx72VRpJKCdokrgFFur6s9zo91wmizACk90WyrnYSW8Qq001c001",
		"CAQACAQOSYsr_FZJmS5YFEWYmRtM614XEfCC4oEay337ggnORRgBL37j_fPgrtd4q001c001",
		"CAQACAQN_xjwSuFXPTzJBLOg6jVs_lfl-WS4vL5Ayj5kVKsMCpQBPjy8SWQBx54kq001c001",
		"CAQACAQN97-Gn595r4X9KAxCc6oAKrj3xBfHDB3fQeUsf0KFLsgDhecOnBeJWB40q001c001",
		"CAQACAQME9AFYJ9RR54f11sVYyS4pCvcrIGZ2cOoDytsMejXl6gDnynZYIPsz-BMq001c001",
		"CAQACAQNIG8OWhBAoCDAUS0yMugUHmI64JA11wmBMLRHEgWX0zgAILXWWJDc4mFsq001c001",
		"CAQACAQMoWZ1_VIXcx0uZPeXXQgpqqLY_76DmFX9sFFtv9woGzADHFOZ_74AY-64q001c001",
		"CAQACAQNRsKTmRjZkhaBn3l29i02NwbWdMYzgN5T0gm2AfGBMvgCFguDmMVPZv4wq001c001",
		"CAQACAQMF-NAO8ovIrAa3lh35iirE7I3VSwnu0I2t-sMoc8e0KACs-u4OS2rVLfcq001c001",
		"CAQACAQMqwA06n2-JlLrG1zj2Ri3QfWmhYumFFMq8v31_d2Os5gCUv4U6YtW0P_0q001c001",
		"CAQACAQNoRBKwLc_P0UZcCf3PcV7FqFYB52RXkuCfNNP3apbKSgDRNFew5zMZ2sIq001c001",
		"CAQACAQNkLowPO-iZQ_UX1AjoeD6ytoJIe4diAryKnEW6XJ5YpQBDnGIPeyPOsUMq001c001",
		"CAQACAQPjxMI6tUbut1t7S91LVek8gopLAtqbskfB3abr9jerYQC33Zs6AozhoH0q001c001",
		"CAQACAQOGhJ4Gw9HhUamI5ffZ-Xq1nIpHWhYTR27NCnFKbN5K0QBRChMGWg_Mbtkq001c001",
		"CAQACAQNPX7fgw2F0U2QGaJTzhpXBrw1AuKMYSAshYVz3SMdGBABTYRjguBuUyrsq001c001",
		"CAQACAQNywplvBTR3U5HOBQG6DWOIXUWNcc3HvMifD1wg7ZzQDgBTD8dvcSM5Ufcq001c001",
		"CAQACAQPqwIm9P0Vp3uubOf5IMchd1NCOWHL4qCM3H_vXnh8qdADeH_i9WJVxI3Yq001c001",
		"CAQACAQOitmSSDFOqX4NZxvqxGOM6kd8TG8-Hce3AXxR9stm9vQBfX4eSG5-9nPsq001c001",
		"CAQACAQOfoA0bi2PsfFAgxmDmoDoTJXzFeV_qpSdhh-RiclhUcQB8h-obecQxulUq001c001",
		"CAQACAQMPk42qx_Qd3zn9HJpY92XrIzPhEzumNOU2xOJGScfb7gDfxKaqE8C95eMq001c001",
		"CAQACAQPBGVASaYa3717aer9bbblQM30FwCg3_epqQO3EPrCcNgDvQDcSwDfoxzQq001c001",
		"CAQACAQPmbvup7RU1Fi0qXElvWZcDmdsdsV95u0djB1FzJLf6uAAWB3mpsZGiqzYq001c001",
		"CAQACAQN6qyCp8j9m58AVpywjYO5FTjiAtpmX7vGsZpkS2-NXjwDnZpeptuoVRREq001c001",
		"CAQACAQOZgmT18JJo5QintLqqNX4mvyCz3opIrxi_Dx0UXA-8gADlD0j13jLSvWUq001c001",
		"CAQACAQNnmWQF4_NL6euS94xhRr_WLh4tUEM7SBx3C3976M3ZzgDpCzsFUEOcps8q001c001",
		"CAQACAQMCbZeyEcuUzQexHL_OaAZtI2zjkFh0wKyrdo_zkCu2aADNdnSykEVM0PMq001c001",
		"CAQACAQPYrVxBGGrj9Sk_UN_d8MdtKN7d4Es4EFW38gyfCG30wQD18jhB4LeZdGwq001c001",
		"CAQACAQMb_AwIq4Cbvr0EP7TxyCSjEXalsAxfc8FtgrGxwCfDNAC-gl8IsLzl9gcq001c001",
		"CAQACAQPHU0LKtXpU-K8UhtjeSp7GG1VOoafyC6TLWbtLe2p0HAD4WfLKoRiKDNMq001c001",
		"CAQACAQOpb22Y_uqY5RQhSAIF7VgjbVXfjOmmOzyh8Hw-IrS8PgDl8KaYjGoK-UYq001c001",
		"CAQACAQMVlECa-7qEwE8Eh0i_RF44EoMtGIwRKTVy2XoSZct6rgDA2RGaGNwwbnsq001c001",
		"CAQACAQOwinelk908VWV8DKqe8zwDLkNUzsleAchSiyWPEPR_LwBVi16lzvh-6TYq001c001",
		"CAQACAQP8zVivPKKeewGo0fCojhLTAURyWh_qEjqZipZf_w5swgB7iuqvWtDTgDgq001c001",
		"CAQACAQNlKZ_GQvhyMyiPfQa2bGRH3_h32QDAjquj0wuLCUNdfgAz08DG2Wh6BcIq001c001",
		"CAQACAQP1tyxteGxd4Xl5o2BXG_XVUqwYfmbEm82nMqkPpgEXVgDhMsRtfphD7j0q001c001",
		"CAQACAQM5K_y3kStlQvI4myIRtHCWdxyXSjgEEVSc-ayzWMKMwgBC-QS3Sh-wJpEq001c001",
		"CAQACAQPhnlw5xnvh7DzCB1O7rxnpwM2AVr-3pnnf3x74vVyqvwDs37c5Vn7sulYq001c001",
		"CAQACAQPAXeVr2gEy67xuYTjXzLSRpQsdeNCyyK3Fv1tShLFw8gDrv7JreE3teaUq001c001",
		"CAQACAQP1Hn1tOB6K8o30rOUzZMrnf5J-jo91n6K7y4OxOFnQlQDyy3VtjiLcMVcq001c001",
		"CAQACAQNKOiT--ow9yORjPSZBompEsyftjAOMH4Uqls_dgIX4kgDIloz-jKsBZKMq001c001",
		"CAQACAQNxwlyov_XbAr-f9Yg2XaII251Fl-HMAOOV1s2dre5EdAAC1syol5fbf30q001c001",
		"CAQACAQOcUQUXhK15ZE97V8iOFACKlXchHA4QebwqJhMZJtm5ngBkJhAXHAFqBvMq001c001",
		"CAQACAQMRqsRJA0KUXbsgMwP6neXjSGh0XqqO2ajzTGhL7XeY0wBdTI5JXlNTMOQq001c001",
		"CAQACAQNfPhW83mnOCZyIO_i7EcWKEcyIy924Pb847nwJm6yJ7wAJ7ri8y74Mwhsq001c001",
		"CAQACAQPUFNly5GPWWtjQV4T5SyyjVkQIYG4dwrsi35heNK6-zwBa3x1yYNtAaXMq001c001",
		"CAQACAQM3MciDL66r3d5-JbFfyJ-XMlHPqU7OnEsI1bsZZ0GHOQDd1c6DqUWi0SUq001c001",
		"CAQACAQPsNFwdSYFH7BrLSdG24VFiqZrLyiHva1-M6W24ChDG9ADs6e8dypmWFkEq001c001",
		"CAQACAQPk6OwA1xgl6QxoMDsZ02EtLemjBHnsX2Qmv38m5zA_FwDpv-wABPzueb4q001c001",
		"CAQACAQN4OZ0a-XFvlClEoFjRlHRldQpPTYU4h2WwO9wMn3Bn6gCUOzgaTaHe7oEq001c001",
		"CAQACAQPPUYTVLOWI7xqWI0GY9xmY8G1gNNLZ2lTs-_Ct-KLgYADv-9nVNOLF6WYq001c001",
		"CAQACAQPMJYbpRi0Nhyixn3LY58_3v3Vn1ApFTV3CDUErLB9pXACHDUXp1PYxB3cq001c001",
		"CAQACAQNCGF_9SkLWIEppJUE9upFOGi0MoSoAWP9JplY0wQx84gAgpgD9oQy_o2wq001c001",
		"CAQACAQNcCrNMzI3RlYh4gnmpxmQyS3jWF_jwCGMrEuofoxvCIQCVEvBMFxgrKcIq001c001",
		"CAQACAQPfeGEkki7jkeSp1lAMWE8sb_1IGfo4TXx6_xZIgxHUBQCR_zgkGRj6pp8q001c001",
		"CAQACAQPyCAKj5ZI43Sz8I1tLqfw01xmqo4V8gtqsNmxh6wCYbQDdNnyjo8srlAYq001c001",
		"CAQACAQNlvAQeBNZBN2JHUfa7N6bL9yTjxw9TBEtsUVnqTl9XIQA3UVMex1JHRKIq001c001",
		"CAQACAQMAY94xwsedX4Ln-JM32Fj5gjuyRudj2gBXgGtgwXyLiQBfgGMxRiwfVFAq001c001",
		"CAQACAQMm-ZNVdyfUEDCXIogFCb-2dz0RGr7tdTQ2LNTXcIxIrQAQLO1VGtYkeggq001c001",
		"CAQACAQMovI4kYNCNXeYwKJgd2K-e497PnjqpAwLl5UJa-M4FdABd5akknsB2lhYq001c001",
		"CAQACAQP34OVfXK0gFscDTq7Dcdjduc1xgk4foJ0GEy63zk7nxQAWEx9fgtDQD8Mq001c001",
		"CAQACAQMvjNJYIJhXbOAlE3xsLhwJYAZb5Mb6EzOOxuXkH9EisgBsxvpY5LATMrgq001c001",
		"CAQACAQNQQbRbEBttEeavhOpNQvPdumUbOTYL2FK-aiu5S_ybzwARagtbOWhShGAq001c001",
		"CAQACAQPjTSaqKjPSj72iZcQcPn2uZRPEMWSu5xJ2_UDkLS6QPwCP_a6qMbzkW2Mq001c001",
		"CAQACAQPiyVp6-oA9Lt6pwzl3SWeLrD5d9OuuD-ZR_ihdLZTDsgAu_q569LplUDYq001c001",
		"CAQACAQO1gON0n9y48wVIgw708arZTT4Kt12QYCGtagg2oUcTxgDzapB0t5qX1Moq001c001",
		"CAQACAQPvDaOAuMXD60OEQZLhjjZZPtpSaFIzZFkhYp-SaT4NlgDrYjOAaHsphFYq001c001",
		"CAQACAQNTnlbXG-ruVgY3b5aTljB-8d0Xa4_a27u8qulNfEmonABWqtrXaxLv2YIq001c001",
		"CAQACAQMlls3eRGAB_YrfwqqLWAwR4v3dNToKdY5X2kHmjundPAD92greNWhqpvoq001c001",
		"CAQACAQOOgmoV3Tv0hWIPVS4uAtrG_h0mWt5Ij_sAfA15n84vkACFfEgVWh-BF2Aq001c001",
		"CAQACAQOZ3O4WjJUxjKhomTad5YvRNU44ll-yf6TrnaluxTErbwCMnbIWlsWteQEq001c001",
		"CAQACAQM0Pi1iTVjMvlKVL_E_pdtgtdbQglVoAYcyfeLbqSUnAAC-fWhigrtgvIcq001c001",
		"CAQACAQMaii3M6Kk_tybli-XEIdWZ38RCwCgGNfnl3OzzsvmFmgC33AbMwALCVG0q001c001",
		"CAQACAQPJ306MRtMB41QHg32HMFSU0rLwxkjaAZ8gY5wPgXeqeADjY9qMxsyISPYq001c001",
		"CAQACAQNu2BRlf44YeAubXSp2wTbcaN7G67Kvrqtv-5Oq3UN0PAB4-69l6zz1eZYq001c001",
		"CAQACAQPKvDV7ch5IIdjTeSlMqHH37jzZJ2qaLUdiAH3uhEBYcQAhAJp7J1S55x4q001c001",
		"CAQACAQPBHVF4hhOz1xrC8gpgoSDAgK-IAqveDZVmS61U6zSykADXS954AthQ63Yq001c001",
		"CAQACAQMgXGCN8tjDi93J4Snq-i_nifkTSu9Yel4uo9rv8j_wUwCLo1iNSu_MWs8q001c001",
		"CAQACAQMs8cgxN9fClmTFeSOg5p4D9EhVSdp_RtNyao8fEc3mJACWan8xST2asR8q001c001",
		"CAQACAQMow8zceK7ZfKR5c3H2igJ7-FMEhedks6vyJWNdw56z8gB8JWTchYYfxfgq001c001",
		"CAQACAQNs77s-7DOJVY-mvUUucwBXfsN81kuRq-f-2_p7F8uYjwBV25E-1ouOJMQq001c001",
		"CAQACAQPZSTpKW4mBAT2H_1NUv1yUc2b2wndDC1OBCZVgkQu_EgABCUNKwjkkGKMq001c001",
		"CAQACAQOeJjYoQ9HbW7k6qi23iFhJ8tpJoopLQTpTDV16AgrolQBbDUsoomOXn_4q001c001",
		"CAQACAQNqOdSNwozhDFc8UgxXe-5DXykhJtrBG7JnRpvDhk3geAAMRsGNJnYdwOwq001c001",
		"CAQACAQNkxII0m4EnToX99KmS5Ol_dm6FVMSd8viCuF1tiiUyYQBOuJ00VNIoUWoq001c001",
		"CAQACAQPdw7Fp9FyPKtBAACc7-q_W18eCg3XQOQ8pMUEJ1bsVzwAqMdBpg9aV3Pcq001c001",
		"CAQACAQN4oeeAG2Zn7wLm1ebwisY6wBkDtu9a1Nru6YzSr8RdmgDv6VqAtmHLRF4q001c001",
		"CAQACAQPsLDw3MwGNFMvQbyBD3T3oPfjGPbP89yubLLgumQzeTwAULPw3Pa-vHf4q001c001",
		"CAQACAQPE5e4WyBSgU-SLis6yoKn-1iU4vSk-nk2a_e8lNv3nbwBT_T4WvdMcJKMq001c001",
		"CAQACAQNU93V4H3IGew3FOONRGhAk9j-Qm5eGZDg4MhwnrpcUjwB7MoZ4m5VAtvwq001c001",
		"CAQACAQN8NnFvrIuWKouFw3k6Bj-Au-rs4017NctDFb_HtRUxLgAqFXtv46NX0WUq001c001",
		"CAQACAQNQFFkaWXyaJ3T-pk-_JwP4Y3VogYZfJk1qrrZUC9N3xgAnrl8agcVm_g0q001c001",
		"CAQACAQP-nUJ-GEcJvPBbqOsPNe_rqN7_UfwWsZJNuWqssJL8FAC8uRZ-UUS7Er4q001c001",
		"CAQACAQNaV_emFGp1GQpKlwk25Y5xHvL49opxKQWFlNa66PcRvgAZlHGm9rrTDOcq001c001",
		"CAQACAQPnC16rSuKsP6aX5ng9pX-ixBeO5elF2Vtq9-oUh2FsSgA_90Wr5Rs8T8oq001c001",
		"CAQACAQMV78olMXz1tR5qAb8NtzDHrHkyKqKrHm38jz_H4oQcXQC1j6slKvi1Dnsq001c001",
		"CAQACAQMnjGWkaOezK7qGX08_FNELFSAys3DtXSrA3Vz2YiUiSgAr3e2ksxESq_Mq001c001",
		"CAQACAQO9-Dtmijm_HhG8VNRrLrY1ox7GMZQ6yeWMESgba1JzfwAeETpmMSMkpWIq001c001",
		"CAQACAQPEo8COkZeqjxhflNnK4wKX4Az4qLaqNXvAdn_3O6YTcQCPdqqOqDfWD9wq001c001",
		"CAQACAQM3pWHXp7bUlcf8SHbbVZGzkrBCkncgwe4R2SSoDGgEmQCV2SDXkmNceaEq001c001",
		"CAQACAQNuRs5_-nEJwDD-bJ5rOHknPxPETAnI9imM9SGBReNjKQDA9ch_TN1fNeMq001c001",
		"CAQACAQNigYUW0dXQ1mWYjmHwy_B7PojXxgFeFPbLJbiReOUA2ADWJV4WxrjlXQ8q001c001",
		"CAQACAQP4CuRXbJwzENga7d8pawHlm73gl2bTvyseU4JKGixaWgAQU9NXl8rjwhQq001c001",
		"CAQACAQNUL85sVgzWJVs0UP5Bo5RO0-UiZVptmsEOhwnzYq-iUwAlh21sZVuJ4F0q001c001",
		"CAQACAQNYR_NRotAQub-LJi6szgZyjoYbZefck20Q6gGE57OtSwC56txRZYJxfw8q001c001",
		"CAQACAQNKaEuRufcro2fEnSJLqBqmosH01Y1HQkwN9aZ6GZjpvACj9UeR1WZR1yIq001c001",
		"CAQACAQParlaD667KsVWz5fmCBtsm7teMx1qAArtZiBC88CeTBwCxiICDxzt9LO0q001c001",
		"CAQACAQMrcroHfwStS99jplGQcDGIZXtELIwwWCnJo4v3PqYjyABLozAHLO7teNQq001c001",
		"CAQACAQO8EFXwBKWZfDXSbNUsfiwLOBnOCRRdFkW6ew2DLv1e5QB8e13wCSEeTaEq001c001",
		"CAQACAQP5sDuIlMP1t8_jMZjcL8Njb45g8kJidU4n-yDlaItXUgC3-2KI8lPn5aEq001c001",
		"CAQACAQOd0y3oo7RMSyJ0cqzqZNjgP8hoFp4WTRtNz60YynWvxABLzxboFvrBfekq001c001",
		"CAQACAQNXfFR2tcQNiXcGb7aNnqfw6zm1fY9gpKWBTj_QUV5OLQCJTmB2fcjAwpMq001c001",
		"CAQACAQObIvLGMu2jTG1nURue7P4-rcebeMZvbCf3kcw6vJ4LTgBMkW_GeLv-84kq001c001",
		"CAQACAQMKdKxCHB9fSRbwtmtbICvlRJD1z5YvkN7tCKlvhcxiyABJCC9Cz6Wd-4Eq001c001",
		"CAQACAQNKcbYA7Dy_ySePE4rz3ff8b6FmJnAxbkB7zFDJdePSMQDJzDEAJnNshDwq001c001",
		"CAQACAQPFxKe6n8wnwhw3gyApNXptZ3ghyiSCujck34ZEaST0gADC34K6yo3fsXUq001c001",
		"CAQACAQPmU-QwezT5subIYkmHl0faxHAR2G8QpgT57spX9ZvsSgCy7hAw2LceVqIq001c001",
		"CAQACAQPBljJo9bec49uJB5Hsu2ovYw-dxZBKEY6lUTAlRBcXogDjUUpoxaxvywkq001c001",
		"CAQACAQOF9vBgR0NqM1EWW2lsWmBudM6FGTKHoP8IQ37AAExpAwAzQ4dgGaaaaWsq001c001",
		"CAQACAQMqoK-tA36NanboSrDRlT4_r6NZGvVGJ7E52Qon7KQM9gBq2UatGkad3jkq001c001",
		"CAQACAQPWiJrRAR8-dYv0PsI3jr9L7j5yHyCeLz7La_B3IMDJ6AB1a57RH64cv6gq001c001",
		"CAQACAQOH-bMaObdxy7qfszoVXKoDPTcOQFETlN_xVidimVdRcADLVhMaQBxtpg0q001c001",
		"CAQACAQNlX0IONU3uv6wNauJHR-hQ_DVLq68TjB7IDSSmDBibfAC_DRMOq01EWk8q001c001",
		"CAQACAQMikiWK6ZBxzsErfz1X__tlUtPWcjHu6bDnhODN099wswDOhO6Kct3Qny4q001c001",
		"CAQACAQM0fskSSdO9zRTdk1QkodbAShvW8hU3NyXyPktkRmbsBADNPjcS8laP_24q001c001",
		"CAQACAQMZaojRR614wLBNAQViJL4b951R4gKcVUNpzKWppAJ0cQDAzJzR4mNWyP8q001c001",
		"CAQACAQPm5vfHrnOmAr-2xU_RC2eQf9_ZgGGW9Vw_hQB6i-vFXwAChZbHgJ69978q001c001",
		"CAQACAQP8NcbA_I51y3VtO45LrV3fJyInFYk3C62-7kXx4J4dbADL7jfAFaWFwSkq001c001",
		"CAQACAQO33Ud5Jjcl1y9ePWiU4fI-7bBRcRMkl5I-bKodwC3OlgDXbCR5ccntpa8q001c001",
		"CAQACAQOavn5z1-MBSquDvYmYS2uoFW9xX_4qFh2OqG7bYX-hOgBKqCpzX8RTOVYq001c001",
		"CAQACAQMc4O-HAIIcHk7vSMYOhlXLWvdy294UEI7E1Ta4-8qPSwAe1RSH200MbNAq001c001",
		"CAQACAQMgbDK8TrohwBJI1OmJIZmXFWsifumwXJ6e4SFDoZ0UfQDA4bC8fsG8Rl4q001c001",
		"CAQACAQPjVWxZFYkBTPQusP5hvXP-nz6R1gY4-hTVXeQHlikzSwBMXThZ1uK3hbAq001c001",
		"CAQACAQMvGirbdAnFq21Wp0aDh_s6UVQp5GdeW1NlrC6_fHK0lACrrF7b5KwxwVMq001c001",
		"CAQACAQNVEAfKd4PfXZ_1OvaSwx6_p_-d6b06l2zjhFK05TKlaABdhDrK6R2UWmQq001c001",
		"CAQACAQO5JMrh9wHiOaZ5JeBWJUeaGqLA5lU1f1wxvaAkVfPqSAA5vTXh5g78GKIq001c001",
		"CAQACAQPTC8pmbYGt0zRAM18Wr16FQa2u8fPzTVcnGm0PxB6VAwDTGvNm8RuZB5Uq001c001",
		"CAQACAQPIBWZ7cTv0e-AUSDafMYWqwdfVxAQNhXOEh98u9L8GbAB7hw17xKVgeBUq001c001",
		"CAQACAQOCpv-7WQirswjwuEZNONaaLnDKd1JQEjVldi9RuqH4rgCzdlC7d2xnTdUq001c001",
		"CAQACAQMsYYGjiY683DQoJy3ZqpdH-j2cdi28OfCqxNMEnzUTkQDcxLyjdrGPTREq001c001",
		"CAQACAQMyjwVZhN2s-jaajqB84xgsQejXlEn3RZiL3J5hdBObHAD63PdZlFpNRSEq001c001",
		"CAQACAQM6kImmOVL2MN9X7IwK_23zvrImM87O9Mw0mKfjEC6mQQAwmM6mM4EgK3Yq001c001",
		"CAQACAQNT_siebBRk-bUQENK-v9kUA4jFYJoRaTFmz_NqYf43MQD5zxGeYID9LJUq001c001",
		"CAQACAQNKFllgSqYVwhPzTMXQ_NV5kK1CARlJ-lD5q8ORVWbDQwDCq0lgAd3be08q001c001",
		"CAQACAQP39dE2Ap76dIRgfQNIPQueIXEztGmGdzGiRGNUHx7-eAB0RIY2tI5cbisq001c001",
		"CAQACAQMoH1ybKU9qIzeYMuLjxeO9d4BoRvWXwWbMdtyYY1OV4QAjdpebRpcLuVcq001c001",
		"CAQACAQPjrZslwJzZITnqWQOUUjfDH-tft5Kgsiekscp--0EOMgAhsaAlt9OYLmkq001c001",
		"CAQACAQMWNZ2ft8gdCq9cvtScv_PIMu-3uiVV5Y-AFcJFSxO0IgAKFVWfuu3PlUsq001c001",
		"CAQACAQOcXfBipPm-YEC2gBZQ2nXu3FpA3YdQmH1RmbeBbt7l7ABgmVBi3ZIEgX8q001c001",
		"CAQACAQNfSan7-WLrWtbWudgYeTvl-wzHrNMmWL5EdrmXu01jogBadib7rLlsbv8q001c001",
		"CAQACAQO7pgElPPnVO7z_Q6JTQR_QhdwGf503BM49pGahcEJPZAA7pDclf2kxPJkq001c001",
		"CAQACAQMUQSdIf11xWkCmpOUaT4o_fIIVX9o8WjvKRqgN094Q8gBaRjxIX0-vLnwq001c001",
		"CAQACAQMSMPr65AUUhYgNolSj9Guu__67nhDfCRkKGFz7Zvz64gCFGN_6nqgGUKsq001c001",
		"CAQACAQOsjeXKAABEpfJCDCkjIlr1c3RZ4oxjUkfPgW293LK_BwClgWPK4k1ZK1Eq001c001",
		"CAQACAQNvwdCFerg4nCmBjamPdCNbhPaXnh5NSKFtN32lDsDV5gCcN02Fnvd32dAq001c001",
		"CAQACAQPBNdQFzD6vj_BXXRptxKbkAct-VGtjJp0SnzU7R-ZMcwCPn2MFVJ2A_ycq001c001",
		"CAQACAQMWeu3PBpgnYKciG9MsVr65bhGmgf27nexk0H4zbn8NvQBg0LvPgRV7Ga8q001c001",
		"CAQACAQPg9HYp8F05w3Lg6YcqAedfL4kCdJUFG6BRXpoBYG8FGADDXgUpdK7EZBgq001c001",
		"CAQACAQMlQBLigajbCUbWH5uMtwcUyq8hy2_jLj79LVI-3rtw3AAJLePiy0g3-N0q001c001",
		"CAQACAQPG2GEsOrs41NX7p_o0KZUtrshj6JOsh-Xf7Lc4EFdRYgDU7Kws6HjyYbgq001c001",
		"CAQACAQNuICtEFyb7DIM2VrBS_t00onISv3xHL8YVtfHKfkML6QAMtUdEv-uQB2wq001c001",
		"CAQACAQM7YgBZLms0-GL-Ne5yeyI20P3Ik2ByGjiNfHO8_AbP5AD4fHJZkzw33RMq001c001",
		"CAQACAQNd-6Xcg7Nc3rOF5XGpBTDfnJD6uaSspGk-6qF4xFRs3wDe6qzcuUZ-w9Mq001c001",
		"CAQACAQPH_44aN9X-qUCRy6_kl4fWLJ7mSOB3vnUSpzzek38fJQCpp3caSI1R9toq001c001",
		"CAQACAQPwo88QD6o2iOUOjYoqSemvTVpt_JpV8kB5SPDmqDmOsQCISFUQ_N-7W20q001c001",
		"CAQACAQPL06QbXxFDXX7ovOx0yCK2OnaK1t1X3j7K0ysOx3N7kQBd01cb1sGfs6Qq001c001",
		"CAQACAQMAUtCSZN4I17f9iVa9A0jICbur1Yj8CoogR4N4ciKDAgDXR_yS1UANArIq001c001",
		"CAQACAQOTR3NOAjGnh4pGoyNs6CoQZKy9j1YDAjOuRgL0v1dDagCHRgNOjwKT5boq001c001",
		"CAQACAQMZitUumrw-81sGLwo9L6wbXZ4pzUrYIRJWl3DtVnI52QDzl9guzdNCiYYq001c001",
		"CAQACAQM_mpA0SZsk7DuwWvEKFU7geei2TUYaeXv9kRUbPV9DSwDskRo0TRzPPo8q001c001",
		"CAQACAQMMUxOlREboAuiTV3erIfKHpHAXpiH7fjDafWXduEcOlQACffulph2Vzrgq001c001",
		"CAQACAQMBuNcZA8hWAOi_i5kt021eTIlZSioakPtyAU4CTRD81AAAARoZSj5HJvwq001c001",
		"CAQACAQMMGLe1he8EbYlOlQfENpIeehZtpz_tJlVXiaK-o3K3egBtie21p4nOkLIq001c001",
		"CAQACAQOXtkMYVmF3CaALx6K7JMvjqleuE2i4PaATGOvwNerh7AAJGLgYEyYcUhcq001c001",
		"CAQACAQPnlQKsbwpTaSixvA0SkhEC26x-AP12NESACYjQIN8SSABpCXasAIVrGMwq001c001",
		"CAQACAQMDZ1u98P8Lyi4H0AN0D9gBfeNiU6EQYqdpHGc75-PO6QDKHBC9UyYh65Yq001c001",
		"CAQACAQPFGaLhbLTSRpc5VRJtfC86tpRyhxpdP2O7sUvZDIsGgQBGsV3hh6AfCGEq001c001",
		"CAQACAQP8ZDRYsh0H1tqXZ3KiJR5OlBekGdlWNwiBXPjh7moCkADWXFZYGYpOlw4q001c001",
	}
)

const (
	baseURL = "http://production-api.acq-visa-clearing.melifrontends.com"
	// baseURL    = "https://internal-api.mercadopago.com/acq/staging/visa-clearing"
	xAuthToken = "eeb8907e83fec666b5573658cea8e856e73da89c3ebeb8490e395cf3af890036"
)

type (
	ARN struct {
		ID    string `json:"id"`
		State string `json:"state"`
	}

	ClearingRequest struct {
		IsSafetyNet bool `json:"is_safety_net"`
	}

	Order struct {
		ID              string          `json:"id"`
		State           string          `json:"state"`
		ARN             string          `json:"arn"`
		TCRecord        []string        `json:"tc_record"`
		SettlementDates SettlementDates `json:"settlement_dates"`
		Revision        Revision        `json:"revision"`
		BatchKey        BatchKey        `json:"batch_key"`
		IsSafetyNet     bool            `json:"is_safety_net"`
		ClearingRequest ClearingRequest `json:"clearing_request"`
	}

	SettlementDates struct {
		Reconciliation string `json:"reconciliation"`
		Settlement     string `json:"settlement"`
		Value          string `json:"value"`
		Merchant       string `json:"merchant"`
		WorkingDays    int    `json:"working_days"`
		CalendarDays   int    `json:"calendar_days"`
		ValidToUTC     string `json:"valid_to_utc"`
	}

	Revision struct {
		UpdatedAt string `json:"updated_at"`
	}

	BatchKey struct {
		UniqueFileID string `json:"unique_file_id"`
		BatchNumber  int    `json:"batch_number"`
	}

	GetStateResponse struct {
		ID    string `json:"transaction_id"`
		State string `json:"state"`
	}

	Authtrx struct {
		ID         string    `json:"transaction_id"`
		ModifiedAt time.Time `json:"modified_at"`
		State      string    `json:"state"`
	}
)

func main() {
	// fmt.Printf("Reprocessing %d orders\n", len(orderIDList))

	var count = 0

	wg := sync.WaitGroup{}
	for _, value := range orderIDList {
		wg.Add(1)
		go executeReviewAndReprocess(&wg, value)

		count++
		if count%49 == 0 {
			time.Sleep(1 * time.Second)
		}
	}

	wg.Wait()
	fmt.Printf("%d order(s) reprocessed", count)
}

func executeReviewAndReprocess(wg *sync.WaitGroup, value string) {
	defer wg.Done()

	order := getByOrderID(value)
	fmt.Printf(
		"ID: %s, St: %s, ARN: %s, Update: %s, FileID: %s, BatchNum: %d\n",
		value,
		order.State,
		order.ARN,
		order.Revision.UpdatedAt,
		order.BatchKey.UniqueFileID,
		order.BatchKey.BatchNumber,
	)

	// order := getByARN(value)

	// fmt.Printf("ARN: %s, St: %s, ID: %s\n", value, order.State, order.ID)

	// if order.State == "presenting" || order.State == "presentation_in_review" || order.State == "sent_exception" {
	// 	if order.State == "presenting" || order.State == "sent_exception" {
	// 		err := reviewOrder(order.ID)
	// 		if err != nil {
	// 			fmt.Println("Error review: %w", err)
	// 			return
	// 		}

	// 		err = reprocessOrder(order.ID)
	// 		if err != nil {
	// 			fmt.Println("Error reprocess: %w", err)
	// 			return
	// 		}
	// 	} else if order.State == "presentation_in_review" {
	// 		err := reprocessOrder(order.ID)
	// 		if err != nil {
	// 			fmt.Println("Error reprocess: %w", err)
	// 			return
	// 		}
	// 	}
	// }
}

func reviewOrder(orderID string) error {
	reviewURL := fmt.Sprintf("%s/maintenance/clearing/orders/%s/review", baseURL, orderID)

	err := doPost(reviewURL)
	if err != nil {
		return err
	}

	fmt.Printf("Order in review: %s\n", orderID)

	return nil
}

func reprocessOrder(orderID string) error {
	reprocessURL := fmt.Sprintf("%s/maintenance/clearing/orders/%s/reprocess", baseURL, orderID)

	err := doPost(reprocessURL)
	if err != nil {
		return err
	}

	fmt.Printf("Order in reprocess: %s\n", orderID)

	return nil
}

func getByOrderID(orderID string) Order {
	getByOrderID := fmt.Sprintf("%s/maintenance/clearing/orders/%s", baseURL, orderID)

	returnedData := doGet(getByOrderID)

	var result Order
	_ = json.Unmarshal(returnedData, &result)

	return result
}

func getByARN(arn string) ARN {
	getOrderByARN := fmt.Sprintf("%s/maintenance/clearing/orders?arn=%s", baseURL, arn)

	returnedData := doGet(getOrderByARN)

	var result ARN
	_ = json.Unmarshal(returnedData, &result)

	return result
}

func doPost(postURL string) error {
	reqURL, _ := url.Parse(postURL)

	req := &http.Request{
		Method: "PUT",
		URL:    reqURL,
		Header: map[string][]string{
			"Content-Type": {"application/json; charset=UTF-8"},
			"x-auth-token": {xAuthToken},
		},
	}

	res, err := http.DefaultClient.Do(req)

	res.Body.Close()

	return err
}

func doGet(getUrl string) []byte {
	reqURL, _ := url.Parse(getUrl)

	req := &http.Request{
		Method: "GET",
		URL:    reqURL,
		Header: map[string][]string{
			"Content-Type": {"application/json; charset=UTF-8"},
			"x-auth-token": {xAuthToken},
		},
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal("Error:", err)
	}

	data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	return data
}

/*
func main() {
   urlGetByState := "https://internal-api.mercadopago.com/acq/internal/visa-authorization/v1/transactions/state/capture_requested"
   urLGetByID := "https://internal-api.mercadopago.com/acq/internal/visa-authorization/v1/transactions/"

   client := http.DefaultClient

   checkLateTransactions(client, urlGetByState, urLGetByID)
   t := time.Tick(time.Hour * 24)
   for {
      select {
      case <-t:
         checkLateTransactions(client, urlGetByState, urLGetByID)
      }
   }
}
*/

func checkLateTransactions(client *http.Client, urlGetByState string, urLGetByID string) {
	fmt.Println("checking late transactions")
	resp, err := client.Get(urlGetByState)
	if err != nil {
		panic(err)
	}

	resps := []GetStateResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&resps); err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	for _, r := range resps {
		wg.Add(1)
		go checkLateTrx(&wg, client, urLGetByID, r)
	}

	wg.Wait()
	fmt.Println("finished checking late transactions")
}

func checkLateTrx(wg *sync.WaitGroup, client *http.Client, urLGetByID string, r GetStateResponse) {
	defer wg.Done()
	getById, err := client.Get(urLGetByID + r.ID)
	if err != nil {
		panic(err)
	}

	trx := Authtrx{}
	if err := json.NewDecoder(getById.Body).Decode(&trx); err != nil {
		panic(err)
	}

	if trx.ModifiedAt.Before(time.Now().Add(-time.Hour * 18)) {
		fmt.Println(fmt.Sprintf("transaction is late: %s", trx.ID))

		req, _ := http.NewRequest(http.MethodPatch, urLGetByID+trx.ID+"/notify",
			&bytes.Buffer{})

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		req, _ = http.NewRequest(http.MethodPost, "https://hooks.slack.com/services/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",
			bytes.NewBufferString(fmt.Sprintf(`{"text": "transaction was late: %s notify response code: %d"}`, trx.ID, resp.StatusCode)))

		client.Do(req)
	}
}

/*
func main() {
   timeNow := time.Now().Add(time.Hour)
   timeLimit := time.Now().Add(365 * 24 * time.Hour)

   fmt.Printf("time now: %s, time limit: %s\n", timeNow, timeLimit)
   fmt.Println(!timeNow.Before(timeLimit))
}
*/
