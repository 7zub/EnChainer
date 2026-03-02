select ccy currency, count(*) k,
       max(case when rn = 1 then tr end) exchange,
       max(case when rn = 1 then spread end) spread_cur,
       max(spread) spread_max,
       min(spread) spread_min
from (select row_number() over (partition by ccy order by create_date desc) rn, buy_ex || '-' || sell_ex tr, * from trade_task) t
group by ccy
order by spread_max desc;


select t.spread, t.stage, t.status, t.message, o.*, r.response_raw, r.code from trade_task t
left join operation_task o on o.task_id = t.task_id
left join request r on r.req_id = o.req_id
where stage = 'trade'
order by o.id desc;

select ccy, buy_ex, sell_ex, max(spread) sp, count(*) cnt from trade_task where spread > 0.4
group by  ccy, buy_ex, sell_ex
order by sp desc;

select * from trade_task where message like '%Низкий объем%'
order by create_date desc;

select * from request where req_type != 'Book'
order by req_date desc;


select p.market,o.* from order_book o
                             join trade_pair p on p.id = o.tp_id
where task_id = '0000665-9771';

select spread, buy_deep, sell_deep, tp.ccy, o.task_id, tt.buy_ex, tt.sell_ex, o.exchange, o.bids, o.asks, o.create_date, o.req_id, tt.message
from order_book o
join trade_pair tp on tp.id = o.tp_id
left join trade_task tt on tt.task_id = o.task_id
left join request r on r.req_id = o.req_id
where tt.task_id = '0006679-6901' --tp.ccy = 'PLAY'
order by o.task_id desc, o.create_date desc;


--UPDATE ex.trade_pair SET sess_time = sess_time + (floor(random() * 1500) + 1) * 1000000 WHERE title = 'F-s';
