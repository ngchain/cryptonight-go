# CrptoNight-go

A new variant of the origin cryptonight algorithm

## Goal?

Fully written on Golang. Can be easily compiled to every platform.

For faster hashrate. The hash faster, the synchronization faster.

For ASIC-resistance. ASIC-resistance is required to keep network safe.

## What's changed?

We studied other CryptoNight variants like lite, turtle upx and so on.

At first, changing the **Iteration** from 524,288 to **65,536**, it will lead to much less iterating simple calculations. The hashrate will get quite a increases.

Then is turn to change the **Scratchpad**. As everyone known, the key to ASIC-resistant in CryotoNote series coins is that Scratchpad make the CryptoNight algorithm memory-hard. We’re dropping the scratch pad size in fourth, from From 2MB to **512KB**. And the hashrate gets around 8 times.

After boosting the speed of hashing, we changed the finalHash period to ensure the ASIC-resistance.

We added 8 more algorithm, and let them randomly occur according by the number, which equals the first byte of the last sum result mods 12. It's similar to X16r, but not same. 

As for why the memory-hard in CryptoNight and randomness in X16r lead to ASIC-resistance, there are more authoritative papers on internet or bitcointalk. So we will not talk this on such a little repo README.

## Ref

most codes forked from `github.com/Equim-chan/cryptonight`.

Modified by ngin-network Team.
