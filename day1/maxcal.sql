--() { :; }; exec psql -f "$0" -d aoc2022_1

-- you need to "createdb aoc2022_1" first if this doesn't start

drop table if exists input;
drop table if exists elfinput;

create table input ( amount varchar )
;

\copy input from 'input/test.txt'

with cte_elfinput as (
    select amount,
           1 + count(*) filter (where amount = '') over ( rows between unbounded preceding and current row ) as elfnum
    from input
)
select amount::int, elfnum
into table elfinput
from cte_elfinput
where amount != ''
;

-- select * from elfinput;

select sum(amount) as total_amount, elfnum
into temp table tmp_topelf
  from elfinput
  group by elfnum
  order by 1 desc
  limit 3
;

\echo top 3 elfs, amount carried and elf number
select * from tmp_topelf;

\echo top 3 together
select sum(total_amount) from tmp_topelf;
