drop account if exists acc1;
create account acc1 admin_name = 'r1' identified by '111' open comment 'acc1';

use mo_catalog;
drop table if exists a;

-- cluster table
CREATE CLUSTER TABLE a (`id` varchar(128) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `account_name` varchar(128) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uniq_acc` (`account_name`));

insert into a accounts(acc1) values (0,"abc","acc1");
insert into a accounts(acc1) values (1,"bcd","acc2");

-- check it in the non-sys account
-- @session:id=2&user=acc1:r1&password=111
use mo_catalog;
select * from a;
-- @session

update a set `account_name` = "cde" where id = 1;
update a set `account_name` = "cde" where id = 1;
update a set `name` = "xxx" where `account_name` = "cde";
update a set `id` = 3 where `id` = 1;
update a set `id` = 3, `account_name`='qwe' where `id` = 3;
delete from a where `id` = 1;

-- check it in the non-sys account
-- @session:id=2&user=acc1:r1&password=111
use mo_catalog;
select * from a;
-- @session

drop table if exists a;
drop account if exists acc1;
