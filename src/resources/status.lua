--redis key
local hMapkey=KEYS[1];
local lastTimekey=ARGV[1];
local lastTimeValue=ARGV[2];
local heartbeatIntervalKey=ARGV[3];
local heartbeatIntervalValue=ARGV[4];
local registerInfoKey=ARGV[5];
local registerInfoValue=ARGV[6];
local expire=ARGV[7];

local numberHeartbeatInterval = tonumber(heartbeatIntervalValue)
--该key过期时间
local expireTime=redis.call("TTL", hMapkey);
local numberExpireTime = tonumber(expireTime)

--设置参数
redis.call('HMSET', hMapkey, lastTimekey, lastTimeValue, heartbeatIntervalKey, heartbeatIntervalValue, registerInfoKey, registerInfoValue);
if (  numberHeartbeatInterval > numberExpireTime ) then
    --设置过期时间
    redis.call('EXPIRE', hMapkey, numberHeartbeatInterval);
end

return 0;