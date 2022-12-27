drop account if exists tenant_test;
create account tenant_test admin_name = 'root' identified by '111' open comment 'tenant_test';

use mo_catalog;
drop table if exists a;
create cluster table a(a int);
insert into a accounts(sys,tenant_test) values(0),(1),(2),(3);
select a from a;

-- @session:id=2&user=tenant_test:root&password=111
use mo_catalog;
select a from a;
-- @session

drop account if exists tenant_test;

select a from a;
drop table if exists a;

drop account if exists tenant_test;