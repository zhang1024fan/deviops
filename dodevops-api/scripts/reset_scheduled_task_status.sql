-- ============================================
-- 定时任务状态重置脚本
-- 用于修复子任务状态卡在运行中(2)的问题
-- ============================================

-- 1. 查看当前所有定时任务及其子任务的状态
SELECT
    tj.id AS task_id,
    tj.name AS task_name,
    tj.type AS task_type,
    tj.status AS parent_status,
    tj.cron_expr,
    tw.id AS work_id,
    tw.template_id,
    tw.status AS work_status,
    CASE tw.status
        WHEN 1 THEN '等待中'
        WHEN 2 THEN '运行中'
        WHEN 3 THEN '成功'
        WHEN 4 THEN '异常'
        ELSE '未知'
    END AS status_name,
    tw.start_time,
    tw.end_time
FROM task_job tj
LEFT JOIN task_work tw ON tj.id = tw.task_id
WHERE tj.type = 2  -- 只查看定时任务
ORDER BY tj.id, tw.id;

-- 2. 重置所有运行中的子任务状态为等待中
UPDATE task_work
SET status = 1
WHERE task_id IN (
    SELECT id FROM task_job WHERE type = 2
)
AND status = 2;

-- 3. 重置所有已完成或异常的子任务状态为等待中（可选）
-- 如果你想让已完成的任务也能重新执行，取消下面的注释
-- UPDATE task_work
-- SET status = 1
-- WHERE task_id IN (
--     SELECT id FROM task_job WHERE type = 2
-- )
-- AND status IN (3, 4);

-- 4. 验证重置后的状态
SELECT
    tj.id AS task_id,
    tj.name AS task_name,
    tj.status AS parent_status,
    COUNT(tw.id) AS total_works,
    SUM(CASE WHEN tw.status = 1 THEN 1 ELSE 0 END) AS pending_count,
    SUM(CASE WHEN tw.status = 2 THEN 1 ELSE 0 END) AS running_count,
    SUM(CASE WHEN tw.status = 3 THEN 1 ELSE 0 END) AS success_count,
    SUM(CASE WHEN tw.status = 4 THEN 1 ELSE 0 END) AS failed_count
FROM task_job tj
LEFT JOIN task_work tw ON tj.id = tw.task_id
WHERE tj.type = 2
GROUP BY tj.id, tj.name, tj.status
ORDER BY tj.id;

-- 5. 针对特定任务的重置（可选）
-- 如果只想重置某个特定任务，使用下面的SQL
-- 将 53, 54 替换为你的任务ID
-- UPDATE task_work SET status = 1 WHERE task_id IN (53, 54) AND status IN (2, 3, 4);