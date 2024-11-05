DROP TABLE IF EXISTS "system_oauth2_access_token";
CREATE TABLE "system_oauth2_access_token"
(
    "id"            INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    "user_id"       INTEGER                           NOT NULL,
    "user_type"     INTEGER                           NOT NULL,
    "access_token"  TEXT                              NOT NULL,
    "refresh_token" TEXT                              NOT NULL,
    "client_id"     TEXT                              NOT NULL,
    "scopes"        TEXT,                                                 -- 授权范围
    "expires_time"  INTEGER                           NOT NULL,           -- 过期时间
    "creator"       INTEGER,                                              -- 创建者
    "create_time"   INTEGER,
    "updater"       INTEGER,                                              -- 更新者
    "update_time"   INTEGER,
    "deleted"       INTEGER                           NOT NULL DEFAULT 0, -- 是否删除，使用 INTEGER 代替
    "tenant_id"     INTEGER                           NOT NULL DEFAULT 0  -- 租户编号
);

DROP TABLE IF EXISTS "system_users";
CREATE TABLE "system_users"
(
    "id"          INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    "username"    TEXT                              NOT NULL,
    "password"    TEXT                              NOT NULL DEFAULT '',
    "nickname"    TEXT                              NOT NULL,
    "remark"      TEXT,
    "dept_id"     INTEGER,
    "post_ids"    TEXT,
    "email"       TEXT                                       DEFAULT '',
    "mobile"      TEXT                                       DEFAULT '',
    "sex"         INTEGER                                    DEFAULT 0,
    "avatar"      TEXT                                       DEFAULT '',
    "status"      INTEGER                           NOT NULL DEFAULT 0,
    "login_ip"    TEXT                                       DEFAULT '',
    "employee_id" TEXT                                       DEFAULT '',
    "login_date"  INTEGER,
    "creator"     INTEGER,
    "create_time" INTEGER,
    "updater"     INTEGER,
    "update_time" INTEGER,
    "deleted"     INTEGER                           NOT NULL DEFAULT 0,
    "tenant_id"   INTEGER                           NOT NULL DEFAULT 0,
    UNIQUE ("username", "tenant_id")
);
INSERT INTO system_users (id, username, password, nickname, remark, dept_id, post_ids, email, mobile, sex, avatar,
                          status, login_ip, employee_id, login_date, creator, create_time, updater, update_time,
                          deleted, tenant_id)
VALUES (1, 'admin', '09936307a478c138a546001365971cded2eed3de6ca5f0e9d4f047d6a74f1808', 'admin', '1', 1, '1', '', '', 0,
        '', 0, '::1', '', 1718616817, '', 1718097781, '1', 1718616817, 0, 0);


DROP TABLE IF EXISTS "system_role";
CREATE TABLE system_role
(
    id                  INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name                TEXT                              NOT NULL,
    code                TEXT                              NOT NULL,
    sort                INTEGER                           NOT NULL,
    data_scope          INTEGER                           NOT NULL DEFAULT 1,-- （1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限）
    data_scope_dept_ids TEXT                              NOT NULL DEFAULT '',
    status              INTEGER                           NOT NULL,
    type                INTEGER                           NOT NULL,
    remark              TEXT,
    creator             INTEGER,
    create_time         INTEGER,
    updater             INTEGER,
    update_time         INTEGER,
    deleted             INTEGER                           NOT NULL DEFAULT 0,
    tenant_id           INTEGER                           NOT NULL DEFAULT 0
);

DROP TABLE IF EXISTS system_user_role;
CREATE TABLE system_user_role
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,           -- 自增编号
    user_id     INTEGER                           NOT NULL,           -- 用户ID
    role_id     INTEGER                           NOT NULL,           -- 角色ID
    creator     INTEGER,                                              -- 创建者
    create_time INTEGER,                                              -- 创建时间
    updater     INTEGER,                                              -- 更新者
    update_time INTEGER,                                              -- 更新时间（需要在应用层设置或使用触发器）
    deleted     INTEGER                           NOT NULL DEFAULT 0, -- 是否删除（0 未删除，1 已删除）
    tenant_id   INTEGER                           NOT NULL DEFAULT 0, -- 租户编号
    UNIQUE ("user_id", "role_id", "tenant_id")
);

DROP TABLE IF EXISTS system_role_menu;
CREATE TABLE system_role_menu
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    role_id     INTEGER NOT NULL,
    menu_id     INTEGER NOT NULL,
    creator     INTEGER,
    create_time INTEGER,
    updater     INTEGER,
    update_time INTEGER,
    deleted     INTEGER NOT NULL DEFAULT 0,
    tenant_id   INTEGER NOT NULL DEFAULT 0
);

DROP TABLE IF EXISTS system_menu;
CREATE TABLE system_menu
(
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    name           TEXT    NOT NULL,
    permission     TEXT    NOT NULL DEFAULT '',
    type           INTEGER NOT NULL,
    sort           INTEGER NOT NULL DEFAULT 0,
    parent_id      INTEGER NOT NULL DEFAULT 0,
    path           TEXT             DEFAULT '',
    icon           TEXT             DEFAULT '#',
    component      TEXT,
    status         INTEGER NOT NULL DEFAULT 0,
    visible        INTEGER NOT NULL DEFAULT 0,
    keep_alive     INTEGER NOT NULL DEFAULT 0,
    application_id INTEGER,
    creator        INTEGER,
    create_time    INTEGER,
    updater        INTEGER,
    update_time    INTEGER,
    deleted        INTEGER NOT NULL DEFAULT 0
);

DROP TABLE IF EXISTS system_post;
CREATE TABLE system_post
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    code        TEXT                              NOT NULL,
    name        TEXT                              NOT NULL,
    sort        INTEGER                           NOT NULL,
    status      INTEGER                           NOT NULL,
    remark      TEXT,
    creator     INTEGER,
    create_time INTEGER,
    updater     INTEGER,
    update_time INTEGER,
    deleted     INTEGER DEFAULT 0                 NOT NULL,
    tenant_id   INTEGER DEFAULT 0                 NOT NULL
);

DROP TABLE IF EXISTS system_dept;
CREATE TABLE system_dept
(
    id             INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,            -- 部门id
    name           TEXT                              NOT NULL DEFAULT '', -- 部门名称
    parent_id      INTEGER                           NOT NULL DEFAULT 0,  -- 父部门id
    sort           INTEGER                           NOT NULL DEFAULT 0,  -- 显示顺序
    leader_user_id INTEGER,                                               -- 负责人
    phone          TEXT,                                                  -- 联系电话
    email          TEXT,                                                  -- 邮箱
    status         INTEGER                           NOT NULL DEFAULT 0,  -- 部门状态（0正常 1停用）
    creator        INTEGER,                                               -- 创建者
    create_time    INTEGER,                                               -- 创建时间（需要在应用层设置）
    updater        INTEGER,                                               -- 更新者
    update_time    INTEGER,                                               -- 更新时间（需要在应用层设置）
    deleted        INTEGER                           NOT NULL DEFAULT 0,  -- 是否删除，使用 0 和 1 替代 bit
    tenant_id      INTEGER                           NOT NULL DEFAULT 0,  -- 租户编号
    unique (name, tenant_id)
);

DROP TABLE IF EXISTS system_user_post;
CREATE TABLE system_user_post
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,           -- id
    user_id     INTEGER                           NOT NULL DEFAULT 0, -- 用户ID
    post_id     INTEGER                           NOT NULL DEFAULT 0, -- 岗位ID
    creator     INTEGER,                                              -- 创建者
    create_time INTEGER,                                              -- 创建时间
    updater     INTEGER,                                              -- 更新者
    update_time INTEGER,                                              -- 更新时间
    deleted     INTEGER                           NOT NULL DEFAULT 0, -- 是否删除  0表示未删除，1表示已删除
    tenant_id   INTEGER                           NOT NULL DEFAULT 0  -- 租户编号
);

DROP TABLE IF EXISTS device;
CREATE TABLE `device`
(
    `id`                  INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,              -- '主键id',
    `device_code`         TEXT                              NOT NULL,              -- '设备编号',
    `guid`                TEXT                                       DEFAULT NULL, -- '设备guid',
    `device_name`         TEXT                              NOT NULL,              -- '设备名称',
    `device_classify_id`  TEXT                                       DEFAULT NULL, -- '设备分类编号',
    `device_brand_id`     TEXT                                       DEFAULT NULL,-- '设备品牌',
    `dept_id`             INTEGER                                    DEFAULT NULL, -- '所在部门',
    `principal_id`        INTEGER                                    DEFAULT NULL, -- '负责人',
    `device_model_id`     TEXT                                       default NULL,
    `thing_model_id`      INTEGER                                    DEFAULT NULL, -- '产品模型id'
    `thing_model_name`    TEXT                                       DEFAULT NULL, -- '产品模型名称'
    `token`               TEXT                                       DEFAULT NULL, -- 'token'
    `asset_status`        INTEGER                                    DEFAULT NULL, -- '设备资产状态',
    `device_status`       INTEGER                                    DEFAULT NULL, -- '设备状态',
    `method`              INTEGER                                    DEFAULT NULL, -- '接入方式',
    `access`              INTEGER                                    DEFAULT NULL, -- '物模接入',
    `image`               TEXT,                                                    -- '图片地址',
    `remark`              TEXT                                       DEFAULT NULL, -- '描述',
    `supplier`            TEXT                                       DEFAULT NULL, -- '供应商厂家',
    `manufacturer`        TEXT                                       DEFAULT NULL, -- '制造商',
    `identification_code` TEXT                                       DEFAULT NULL, -- '二维码识别码',
    `money`               TEXT                                       DEFAULT NULL, -- '金额',
    `enter_time`          INTEGER                                    DEFAULT NULL, -- '入厂时间',
    `production_time`     INTEGER                                    DEFAULT NULL, -- '生产时间',
    `purchase_time`       INTEGER                                    DEFAULT NULL, -- '购买时间',
    `report_time`         INTEGER                                    DEFAULT NULL, -- '保修时间',
    `install_time`        INTEGER                                    DEFAULT NULL, -- '安装时间',
    `use_time`            INTEGER                                    DEFAULT NULL, -- '使用时间',
    `position`            TEXT                                       DEFAULT NULL, -- '安装位置',
    `creator`             INTEGER,                                                 -- '创建者',
    `create_time`         INTEGER,                                                 -- '创建时间',
    `updater`             INTEGER,                                                 -- '更新者',
    `update_time`         INTEGER,                                                 -- '更新时间',
    `deleted`             INTEGER                           NOT NULL DEFAULT 0,    -- '是否删除',  -- 使用INTEGER代替bit(1)
    `tenant_id`           INTEGER                           NOT NULL DEFAULT 0,    -- '租户编号',
    unique (device_code, tenant_id)
);

DROP TABLE IF EXISTS device_brand;
CREATE TABLE `device_brand`
(
    `id`                INTEGER PRIMARY KEY AUTOINCREMENT,
    `device_brand_code` TEXT    NOT NULL,
    `device_brand_name` TEXT    NOT NULL,
    `is_enable`         INTEGER NOT NULL,
    `remark`            TEXT,
    `creator`           INTEGER,
    `create_time`       INTEGER,
    `updater`           INTEGER,
    `update_time`       INTEGER,
    `deleted`           INTEGER DEFAULT 0,
    `tenant_id`         INTEGER DEFAULT 0
);

DROP TABLE IF EXISTS `device_classify`;
CREATE TABLE `device_classify`
(
    `id`                   INTEGER PRIMARY KEY AUTOINCREMENT,
    `parent_id`            INTEGER          DEFAULT 0,
    `is_system`            INTEGER          DEFAULT 0,
    `ancestor`             TEXT             DEFAULT '0',
    `sort`                 INTEGER          DEFAULT 0,
    `device_classify_code` TEXT    NOT NULL,
    `device_classify_name` TEXT    NOT NULL,
    `is_enable`            INTEGER          DEFAULT 0,
    `remark`               TEXT,
    `creator`              INTEGER,
    `create_time`          INTEGER,
    `updater`              INTEGER,
    `update_time`          INTEGER,
    `deleted`              INTEGER          DEFAULT 0,
    `tenant_id`            INTEGER NOT NULL DEFAULT 0
);

DROP TABLE IF EXISTS device_model;
CREATE TABLE device_model
(
    id                INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,          -- 主键id
    device_type_id    INTEGER,                                             -- 设备分类id
    device_model_code TEXT                              NOT NULL,          -- 规格编号
    device_model_name TEXT                              NOT NULL,          -- 规格名称
    version           TEXT                              NOT NULL,          -- 版本号
    is_enable         INTEGER                           NOT NULL,          -- 是否启用
    remark            TEXT,                                                -- 备注
    creator           INTEGER,                                             -- 创建者
    create_time       INTEGER,                                             -- 创建时间
    updater           INTEGER,                                             -- 更新者
    update_time       INTEGER,                                             -- 更新时间
    deleted           INTEGER                                    DEFAULT 0,
    tenant_id         INTEGER                           NOT NULL DEFAULT 0 -- 租户编号
);

DROP TABLE IF EXISTS gateway_classify;
CREATE TABLE gateway_classify
(
    id                    INTEGER PRIMARY KEY AUTOINCREMENT,
    gateway_classify_name TEXT    NOT NULL,
    is_enable             INTEGER NOT NULL DEFAULT 0,
    remark                TEXT,
    creator               INTEGER,
    create_time           INTEGER,
    updater               INTEGER,
    update_time           INTEGER,
    deleted               INTEGER NOT NULL DEFAULT 0,
    tenant_id             INTEGER NOT NULL DEFAULT 0
);

DROP TABLE IF EXISTS gateway_brand;
CREATE TABLE gateway_brand
(
    id                 INTEGER PRIMARY KEY AUTOINCREMENT,
    gateway_brand_code TEXT    NOT NULL,
    gateway_brand_name TEXT    NOT NULL,
    is_enable          INTEGER NOT NULL DEFAULT 0,
    remark             TEXT,
    creator            INTEGER,
    create_time        INTEGER,
    updater            INTEGER,
    update_time        INTEGER,
    deleted            INTEGER NOT NULL DEFAULT 0,
    tenant_id          INTEGER NOT NULL DEFAULT 0
);

DROP TABLE IF EXISTS gateway_model;
CREATE TABLE gateway_model
(
    id                        INTEGER PRIMARY KEY AUTOINCREMENT,
    gateway_type_code         TEXT    NOT NULL,
    gateway_type_name         TEXT    NOT NULL,
    is_enable                 INTEGER NOT NULL,
    remark                    TEXT,
    creator                   INTEGER,
    create_time               INTEGER NOT NULL,
    updater                   INTEGER,
    update_time               INTEGER NOT NULL,
    deleted                   INTEGER NOT NULL DEFAULT 0,
    tenant_id                 INTEGER NOT NULL DEFAULT 0,
    os_type                   INTEGER,
    cpu_type                  INTEGER,
    gateway_classification_id INTEGER
);

DROP TABLE IF EXISTS thing_model;
CREATE TABLE thing_model
(
    `id`               INTEGER PRIMARY KEY AUTOINCREMENT,
    `device_type_id`   TEXT             DEFAULT NULL,
    `thing_model_name` TEXT             DEFAULT NULL,
    `method`           INTEGER          DEFAULT NULL,
    `is_enable`        INTEGER          DEFAULT NULL,
    `remark`           TEXT             DEFAULT NULL,
    `access`           INTEGER          DEFAULT NULL,
    `file_urls`        TEXT             DEFAULT NULL,
    `model_file_id`    INTEGER          DEFAULT NULL,
    `creator`          INTEGER,
    `create_time`      INTEGER NOT NULL,
    `updater`          INTEGER,
    `update_time`      INTEGER NOT NULL,
    `deleted`          INTEGER NOT NULL DEFAULT 0,
    `tenant_id`        INTEGER NOT NULL DEFAULT 0
);

DROP TABLE IF EXISTS gateway;
CREATE TABLE gateway
(
    `id`                  INTEGER PRIMARY KEY AUTOINCREMENT,
    `gateway_code`        TEXT    NOT NULL,
    `gateway_name`        TEXT    NOT NULL,
    `gateway_classify_id` INTEGER NOT NULL,
    `gateway_model_id`    INTEGER NOT NULL,
    `gateway_brand_id`    INTEGER NOT NULL,
    `is_enable`           INTEGER NOT NULL,
    `dept_id`             INTEGER NOT NULL,
    `address`             TEXT             DEFAULT NULL,
    `principal_id`        INTEGER          DEFAULT NULL,
    `image`               TEXT             DEFAULT NULL,
    `online_time`         DATETIME         DEFAULT NULL,
    `registration_time`   DATETIME         DEFAULT NULL,
    `identification_code` TEXT             DEFAULT NULL,
    `status`              INTEGER          DEFAULT NULL,
    `remark`              TEXT             DEFAULT NULL,
    `creator`             INTEGER,
    `create_time`         INTEGER NOT NULL,
    `updater`             INTEGER,
    `update_time`         INTEGER NOT NULL,
    `deleted`             INTEGER NOT NULL DEFAULT 0,
    `tenant_id`           INTEGER NOT NULL DEFAULT 0
);

DROP TABLE IF EXISTS `gateway_device`;
CREATE TABLE `gateway_device`
(
    `id`            INTEGER PRIMARY KEY AUTOINCREMENT,
    `gateway_id`    INTEGER NOT NULL,
    `device_id`     INTEGER NOT NULL,
    `net_interface` TEXT             DEFAULT NULL,
    `agreement_id`  INTEGER          DEFAULT NULL,
    `ip`            TEXT             DEFAULT NULL,
    `status`        INTEGER          DEFAULT NULL,
    `issue_time`    INTEGER          DEFAULT NULL,
    `issue`         INTEGER          DEFAULT NULL,
    `creator`       INTEGER,
    `create_time`   INTEGER,
    `updater`       INTEGER,
    `update_time`   INTEGER,
    `deleted`       INTEGER NOT NULL DEFAULT 0,
    `tenant_id`     INTEGER NOT NULL DEFAULT 0,
    `version`       INTEGER NOT NULL DEFAULT 0
);

DROP TABLE IF EXISTS `agreement`;
CREATE TABLE `agreement`
(
    `id`             INTEGER PRIMARY KEY AUTOINCREMENT,
    `os_type`        INTEGER          DEFAULT NULL,
    `cpu_type`       INTEGER          DEFAULT NULL,
    `agreement_name` TEXT             DEFAULT NULL,
    `file_id`        INTEGER          DEFAULT NULL,
    `remark`         TEXT             DEFAULT NULL,
    `creator`        INTEGER,
    `create_time`    INTEGER,
    `updater`        INTEGER,
    `update_time`    INTEGER,
    `deleted`        INTEGER NOT NULL DEFAULT 0,
    `tenant_id`      INTEGER NOT NULL DEFAULT 0,
    `file_name`      TEXT             DEFAULT NULL
);

DROP TABLE IF EXISTS `gateway_agreement`;
CREATE TABLE `gateway_agreement`
(
    `id`                   INTEGER PRIMARY KEY AUTOINCREMENT,
    `agreement_id`         INTEGER NOT NULL,
    `gateway_id`           INTEGER NOT NULL,
    `creator`              INTEGER,
    `create_time`          INTEGER,
    `updater`              INTEGER,
    `update_time`          INTEGER,
    `deleted`              INTEGER NOT NULL DEFAULT 0,
    `tenant_id`            INTEGER NOT NULL DEFAULT 0,
    `old_agreement_status` INTEGER          DEFAULT NULL,
    `new_agreement_status` INTEGER          DEFAULT NULL,
    `issue_status`         INTEGER          DEFAULT NULL,
    `issue_user_id`        INTEGER          DEFAULT NULL
);