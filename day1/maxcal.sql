--() { :; }; exec psql -f "$0" -d aoc2022_1

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
  from elfinput
  group by elfnum
  order by 1 desc
  limit 1
;
