# Result benchstat

```shell
                 │     old.out      │             new.out             │
                 │      sec/op      │    sec/op     vs base           │
GetDomainStat-16   233402.22µ ± ∞ ¹   20.18µ ± ∞ ¹  ~ (p=1.000 n=1) ²
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                 │      old.out       │             new.out              │
                 │        B/op        │     B/op       vs base           │
GetDomainStat-16   199263.605Ki ± ∞ ¹   7.951Ki ± ∞ ¹  ~ (p=1.000 n=1) ²
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                 │     old.out      │            new.out             │
                 │    allocs/op     │  allocs/op   vs base           │
GetDomainStat-16   1900080.00 ± ∞ ¹   63.00 ± ∞ ¹  ~ (p=1.000 n=1) ²
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
```