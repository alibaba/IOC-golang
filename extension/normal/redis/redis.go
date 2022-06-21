/*
 * Copyright (c) 2022, Alibaba Group;
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis"

	"github.com/alibaba/ioc-golang/autowire"
)

// +ioc:autowire=true
// +ioc:autowire:type=normal
// +ioc:autowire:paramType=Config
// +ioc:autowire:constructFunc=New

type Redis struct {
	client *redis.Client
}

var _ redis.Cmdable = &Redis{}

func (r *Redis) Command() *redis.CommandsInfoCmd {
	return r.client.Command()
}

func (r *Redis) ClientGetName() *redis.StringCmd {
	return r.client.ClientGetName()
}

func (r *Redis) Echo(message interface{}) *redis.StringCmd {
	return r.client.Echo(message)
}

func (r *Redis) Ping() *redis.StatusCmd {
	return r.client.Ping()
}

func (r *Redis) Quit() *redis.StatusCmd {
	return r.client.Quit()
}

func (r *Redis) Del(keys ...string) *redis.IntCmd {
	return r.client.Del(keys...)
}

func (r *Redis) Unlink(keys ...string) *redis.IntCmd {
	return r.client.Unlink(keys...)
}

func (r *Redis) Dump(key string) *redis.StringCmd {
	return r.client.Dump(key)
}

func (r *Redis) Exists(keys ...string) *redis.IntCmd {
	return r.client.Exists(keys...)
}

func (r *Redis) Expire(key string, expiration time.Duration) *redis.BoolCmd {
	return r.client.Expire(key, expiration)
}

func (r *Redis) ExpireAt(key string, tm time.Time) *redis.BoolCmd {
	return r.client.ExpireAt(key, tm)
}

func (r *Redis) Keys(pattern string) *redis.StringSliceCmd {
	return r.client.Keys(pattern)
}

func (r *Redis) Migrate(host, port, key string, db int64, timeout time.Duration) *redis.StatusCmd {
	return r.client.Migrate(host, port, key, db, timeout)
}

func (r *Redis) Move(key string, db int64) *redis.BoolCmd {
	return r.client.Move(key, db)
}

func (r *Redis) ObjectRefCount(key string) *redis.IntCmd {
	return r.client.ObjectRefCount(key)
}

func (r *Redis) ObjectEncoding(key string) *redis.StringCmd {
	return r.client.ObjectEncoding(key)
}

func (r *Redis) ObjectIdleTime(key string) *redis.DurationCmd {
	return r.client.ObjectIdleTime(key)
}

func (r *Redis) Persist(key string) *redis.BoolCmd {
	return r.client.Persist(key)
}

func (r *Redis) PExpire(key string, expiration time.Duration) *redis.BoolCmd {
	return r.client.PExpire(key, expiration)
}

func (r *Redis) PExpireAt(key string, tm time.Time) *redis.BoolCmd {
	return r.client.PExpireAt(key, tm)
}

func (r *Redis) PTTL(key string) *redis.DurationCmd {
	return r.client.PTTL(key)
}

func (r *Redis) RandomKey() *redis.StringCmd {
	return r.client.RandomKey()
}

func (r *Redis) Rename(key, newkey string) *redis.StatusCmd {
	return r.client.Rename(key, newkey)
}

func (r *Redis) RenameNX(key, newkey string) *redis.BoolCmd {
	return r.client.RenameNX(key, newkey)
}

func (r *Redis) Restore(key string, ttl time.Duration, value string) *redis.StatusCmd {
	return r.client.Restore(key, ttl, value)
}

func (r *Redis) RestoreReplace(key string, ttl time.Duration, value string) *redis.StatusCmd {
	return r.client.RestoreReplace(key, ttl, value)
}

func (r *Redis) Sort(key string, sort *redis.Sort) *redis.StringSliceCmd {
	return r.client.Sort(key, sort)
}

func (r *Redis) SortStore(key, store string, sort *redis.Sort) *redis.IntCmd {
	return r.client.SortStore(key, store, sort)
}

func (r *Redis) SortInterfaces(key string, sort *redis.Sort) *redis.SliceCmd {
	return r.client.SortInterfaces(key, sort)
}

func (r *Redis) Touch(keys ...string) *redis.IntCmd {
	return r.client.Touch(keys...)
}

func (r *Redis) TTL(key string) *redis.DurationCmd {
	return r.client.TTL(key)
}

func (r *Redis) Type(key string) *redis.StatusCmd {
	return r.client.Type(key)
}

func (r *Redis) Scan(cursor uint64, match string, count int64) *redis.ScanCmd {
	return r.client.Scan(cursor, match, count)
}

func (r *Redis) SScan(key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return r.client.SScan(key, cursor, match, count)
}

func (r *Redis) HScan(key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return r.client.HScan(key, cursor, match, count)
}

func (r *Redis) ZScan(key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return r.client.ZScan(key, cursor, match, count)
}

func (r *Redis) Append(key, value string) *redis.IntCmd {
	return r.client.Append(key, value)
}

func (r *Redis) BitCount(key string, bitCount *redis.BitCount) *redis.IntCmd {
	return r.client.BitCount(key, bitCount)
}

func (r *Redis) BitOpAnd(destKey string, keys ...string) *redis.IntCmd {
	return r.client.BitOpAnd(destKey, keys...)
}

func (r *Redis) BitOpOr(destKey string, keys ...string) *redis.IntCmd {
	return r.client.BitOpOr(destKey, keys...)
}

func (r *Redis) BitOpXor(destKey string, keys ...string) *redis.IntCmd {
	return r.client.BitOpXor(destKey, keys...)
}

func (r *Redis) BitOpNot(destKey string, key string) *redis.IntCmd {
	return r.client.BitOpNot(destKey, key)
}

func (r *Redis) BitPos(key string, bit int64, pos ...int64) *redis.IntCmd {
	return r.client.BitPos(key, bit, pos...)
}

func (r *Redis) Decr(key string) *redis.IntCmd {
	return r.client.Decr(key)
}

func (r *Redis) DecrBy(key string, decrement int64) *redis.IntCmd {
	return r.client.DecrBy(key, decrement)
}

func (r *Redis) Get(key string) *redis.StringCmd {
	return r.client.Get(key)
}

func (r *Redis) GetBit(key string, offset int64) *redis.IntCmd {
	return r.client.GetBit(key, offset)
}

func (r *Redis) GetRange(key string, start, end int64) *redis.StringCmd {
	return r.client.GetRange(key, start, end)
}

func (r *Redis) GetSet(key string, value interface{}) *redis.StringCmd {
	return r.client.GetSet(key, value)
}

func (r *Redis) Incr(key string) *redis.IntCmd {
	return r.client.Incr(key)
}

func (r *Redis) IncrBy(key string, value int64) *redis.IntCmd {
	return r.client.IncrBy(key, value)
}

func (r *Redis) IncrByFloat(key string, value float64) *redis.FloatCmd {
	return r.client.IncrByFloat(key, value)
}

func (r *Redis) MGet(keys ...string) *redis.SliceCmd {
	return r.client.MGet(keys...)
}

func (r *Redis) MSet(pairs ...interface{}) *redis.StatusCmd {
	return r.client.MSet(pairs...)
}

func (r *Redis) MSetNX(pairs ...interface{}) *redis.BoolCmd {
	return r.client.MSetNX(pairs...)
}

func (r *Redis) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return r.client.Set(key, value, expiration)
}

func (r *Redis) SetBit(key string, offset int64, value int) *redis.IntCmd {
	return r.client.SetBit(key, offset, value)
}

func (r *Redis) SetNX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return r.client.SetNX(key, value, expiration)
}

func (r *Redis) SetXX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return r.client.SetXX(key, value, expiration)
}

func (r *Redis) SetRange(key string, offset int64, value string) *redis.IntCmd {
	return r.client.SetRange(key, offset, value)
}

func (r *Redis) StrLen(key string) *redis.IntCmd {
	return r.client.StrLen(key)
}

func (r *Redis) HDel(key string, fields ...string) *redis.IntCmd {
	return r.client.HDel(key, fields...)
}

func (r *Redis) HExists(key, field string) *redis.BoolCmd {
	return r.client.HExists(key, field)
}

func (r *Redis) HGet(key, field string) *redis.StringCmd {
	return r.client.HGet(key, field)
}

func (r *Redis) HGetAll(key string) *redis.StringStringMapCmd {
	return r.client.HGetAll(key)
}

func (r *Redis) HIncrBy(key, field string, incr int64) *redis.IntCmd {
	return r.client.HIncrBy(key, field, incr)
}

func (r *Redis) HIncrByFloat(key, field string, incr float64) *redis.FloatCmd {
	return r.client.HIncrByFloat(key, field, incr)
}

func (r *Redis) HKeys(key string) *redis.StringSliceCmd {
	return r.client.HKeys(key)
}

func (r *Redis) HLen(key string) *redis.IntCmd {
	return r.client.HLen(key)
}

func (r *Redis) HMGet(key string, fields ...string) *redis.SliceCmd {
	return r.client.HMGet(key, fields...)
}

func (r *Redis) HMSet(key string, fields map[string]interface{}) *redis.StatusCmd {
	return r.client.HMSet(key, fields)
}

func (r *Redis) HSet(key, field string, value interface{}) *redis.BoolCmd {
	return r.client.HSet(key, field, value)
}

func (r *Redis) HSetNX(key, field string, value interface{}) *redis.BoolCmd {
	return r.client.HSetNX(key, field, value)
}

func (r *Redis) HVals(key string) *redis.StringSliceCmd {
	return r.client.HVals(key)
}

func (r *Redis) BLPop(timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	return r.client.BLPop(timeout, keys...)
}

func (r *Redis) BRPop(timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	return r.client.BRPop(timeout, keys...)
}

func (r *Redis) BRPopLPush(source, destination string, timeout time.Duration) *redis.StringCmd {
	return r.client.BRPopLPush(source, destination, timeout)
}

func (r *Redis) LIndex(key string, index int64) *redis.StringCmd {
	return r.client.LIndex(key, index)
}

func (r *Redis) LInsert(key, op string, pivot, value interface{}) *redis.IntCmd {
	return r.client.LInsert(key, op, pivot, value)
}

func (r *Redis) LInsertBefore(key string, pivot, value interface{}) *redis.IntCmd {
	return r.client.LInsertBefore(key, pivot, value)
}

func (r *Redis) LInsertAfter(key string, pivot, value interface{}) *redis.IntCmd {
	return r.client.LInsertAfter(key, pivot, value)
}

func (r *Redis) LLen(key string) *redis.IntCmd {
	return r.client.LLen(key)
}

func (r *Redis) LPop(key string) *redis.StringCmd {
	return r.client.LPop(key)
}

func (r *Redis) LPush(key string, values ...interface{}) *redis.IntCmd {
	return r.client.LPush(key, values...)
}

func (r *Redis) LPushX(key string, value interface{}) *redis.IntCmd {
	return r.client.LPushX(key, value)
}

func (r *Redis) LRange(key string, start, stop int64) *redis.StringSliceCmd {
	return r.client.LRange(key, start, stop)
}

func (r *Redis) LRem(key string, count int64, value interface{}) *redis.IntCmd {
	return r.client.LRem(key, count, value)
}

func (r *Redis) LSet(key string, index int64, value interface{}) *redis.StatusCmd {
	return r.client.LSet(key, index, value)
}

func (r *Redis) LTrim(key string, start, stop int64) *redis.StatusCmd {
	return r.client.LTrim(key, start, stop)
}

func (r *Redis) RPop(key string) *redis.StringCmd {
	return r.client.RPop(key)
}

func (r *Redis) RPopLPush(source, destination string) *redis.StringCmd {
	return r.client.RPopLPush(source, destination)
}

func (r *Redis) RPush(key string, values ...interface{}) *redis.IntCmd {
	return r.client.RPush(key, values...)
}

func (r *Redis) RPushX(key string, value interface{}) *redis.IntCmd {
	return r.client.RPushX(key, value)
}

func (r *Redis) SAdd(key string, members ...interface{}) *redis.IntCmd {
	return r.client.SAdd(key, members...)
}

func (r *Redis) SCard(key string) *redis.IntCmd {
	return r.client.SCard(key)
}

func (r *Redis) SDiff(keys ...string) *redis.StringSliceCmd {
	return r.client.SDiff(keys...)
}

func (r *Redis) SDiffStore(destination string, keys ...string) *redis.IntCmd {
	return r.client.SDiffStore(destination, keys...)
}

func (r *Redis) SInter(keys ...string) *redis.StringSliceCmd {
	return r.client.SInter(keys...)
}

func (r *Redis) SInterStore(destination string, keys ...string) *redis.IntCmd {
	return r.client.SInterStore(destination, keys...)
}

func (r *Redis) SIsMember(key string, member interface{}) *redis.BoolCmd {
	return r.client.SIsMember(key, member)
}

func (r *Redis) SMembers(key string) *redis.StringSliceCmd {
	return r.client.SMembers(key)
}

func (r *Redis) SMembersMap(key string) *redis.StringStructMapCmd {
	return r.client.SMembersMap(key)
}

func (r *Redis) SMove(source, destination string, member interface{}) *redis.BoolCmd {
	return r.client.SMove(source, destination, member)
}

func (r *Redis) SPop(key string) *redis.StringCmd {
	return r.client.SPop(key)
}

func (r *Redis) SPopN(key string, count int64) *redis.StringSliceCmd {
	return r.client.SPopN(key, count)
}

func (r *Redis) SRandMember(key string) *redis.StringCmd {
	return r.client.SRandMember(key)
}

func (r *Redis) SRandMemberN(key string, count int64) *redis.StringSliceCmd {
	return r.client.SRandMemberN(key, count)
}

func (r *Redis) SRem(key string, members ...interface{}) *redis.IntCmd {
	return r.client.SRem(key, members...)
}

func (r *Redis) SUnion(keys ...string) *redis.StringSliceCmd {
	return r.client.SUnion(keys...)
}

func (r *Redis) SUnionStore(destination string, keys ...string) *redis.IntCmd {
	return r.client.SUnionStore(destination, keys...)
}

func (r *Redis) XAdd(a *redis.XAddArgs) *redis.StringCmd {
	return r.client.XAdd(a)
}

func (r *Redis) XDel(stream string, ids ...string) *redis.IntCmd {
	return r.client.XDel(stream, ids...)
}

func (r *Redis) XLen(stream string) *redis.IntCmd {
	return r.client.XLen(stream)
}

func (r *Redis) XRange(stream, start, stop string) *redis.XMessageSliceCmd {
	return r.client.XRange(stream, start, stop)
}

func (r *Redis) XRangeN(stream, start, stop string, count int64) *redis.XMessageSliceCmd {
	return r.client.XRangeN(stream, start, stop, count)
}

func (r *Redis) XRevRange(stream string, start, stop string) *redis.XMessageSliceCmd {
	return r.client.XRevRange(stream, start, stop)
}

func (r *Redis) XRevRangeN(stream string, start, stop string, count int64) *redis.XMessageSliceCmd {
	return r.client.XRevRangeN(stream, start, stop, count)
}

func (r *Redis) XRead(a *redis.XReadArgs) *redis.XStreamSliceCmd {
	return r.client.XRead(a)
}

func (r *Redis) XReadStreams(streams ...string) *redis.XStreamSliceCmd {
	return r.client.XReadStreams(streams...)
}

func (r *Redis) XGroupCreate(stream, group, start string) *redis.StatusCmd {
	return r.client.XGroupCreate(stream, group, start)
}

func (r *Redis) XGroupCreateMkStream(stream, group, start string) *redis.StatusCmd {
	return r.client.XGroupCreateMkStream(stream, group, start)
}

func (r *Redis) XGroupSetID(stream, group, start string) *redis.StatusCmd {
	return r.client.XGroupSetID(stream, group, start)
}

func (r *Redis) XGroupDestroy(stream, group string) *redis.IntCmd {
	return r.client.XGroupDestroy(stream, group)
}

func (r *Redis) XGroupDelConsumer(stream, group, consumer string) *redis.IntCmd {
	return r.client.XGroupDelConsumer(stream, group, consumer)
}

func (r *Redis) XReadGroup(a *redis.XReadGroupArgs) *redis.XStreamSliceCmd {
	return r.client.XReadGroup(a)
}

func (r *Redis) XAck(stream, group string, ids ...string) *redis.IntCmd {
	return r.client.XAck(stream, group, ids...)
}

func (r *Redis) XPending(stream, group string) *redis.XPendingCmd {
	return r.client.XPending(stream, group)
}

func (r *Redis) XPendingExt(a *redis.XPendingExtArgs) *redis.XPendingExtCmd {
	return r.client.XPendingExt(a)
}

func (r *Redis) XClaim(a *redis.XClaimArgs) *redis.XMessageSliceCmd {
	return r.client.XClaim(a)
}

func (r *Redis) XClaimJustID(a *redis.XClaimArgs) *redis.StringSliceCmd {
	return r.client.XClaimJustID(a)
}

func (r *Redis) XTrim(key string, maxLen int64) *redis.IntCmd {
	return r.client.XTrim(key, maxLen)
}

func (r *Redis) XTrimApprox(key string, maxLen int64) *redis.IntCmd {
	return r.client.XTrimApprox(key, maxLen)
}

func (r *Redis) BZPopMax(timeout time.Duration, keys ...string) *redis.ZWithKeyCmd {
	return r.client.BZPopMax(timeout, keys...)
}

func (r *Redis) BZPopMin(timeout time.Duration, keys ...string) *redis.ZWithKeyCmd {
	return r.client.BZPopMin(timeout, keys...)
}

func (r *Redis) ZAdd(key string, members ...redis.Z) *redis.IntCmd {
	return r.client.ZAdd(key, members...)
}

func (r *Redis) ZAddNX(key string, members ...redis.Z) *redis.IntCmd {
	return r.client.ZAddNX(key, members...)
}

func (r *Redis) ZAddXX(key string, members ...redis.Z) *redis.IntCmd {
	return r.client.ZAddXX(key, members...)
}

func (r *Redis) ZAddCh(key string, members ...redis.Z) *redis.IntCmd {
	return r.client.ZAddCh(key, members...)
}

func (r *Redis) ZAddNXCh(key string, members ...redis.Z) *redis.IntCmd {
	return r.client.ZAddNXCh(key, members...)
}

func (r *Redis) ZAddXXCh(key string, members ...redis.Z) *redis.IntCmd {
	return r.client.ZAddXXCh(key, members...)
}

func (r *Redis) ZIncr(key string, member redis.Z) *redis.FloatCmd {
	return r.client.ZIncr(key, member)
}

func (r *Redis) ZIncrNX(key string, member redis.Z) *redis.FloatCmd {
	return r.client.ZIncrNX(key, member)
}

func (r *Redis) ZIncrXX(key string, member redis.Z) *redis.FloatCmd {
	return r.client.ZIncrXX(key, member)
}

func (r *Redis) ZCard(key string) *redis.IntCmd {
	return r.client.ZCard(key)
}

func (r *Redis) ZCount(key, min, max string) *redis.IntCmd {
	return r.client.ZCount(key, min, max)
}

func (r *Redis) ZLexCount(key, min, max string) *redis.IntCmd {
	return r.client.ZLexCount(key, min, max)
}

func (r *Redis) ZIncrBy(key string, increment float64, member string) *redis.FloatCmd {
	return r.client.ZIncrBy(key, increment, member)
}

func (r *Redis) ZInterStore(destination string, store redis.ZStore, keys ...string) *redis.IntCmd {
	return r.client.ZInterStore(destination, store, keys...)
}

func (r *Redis) ZPopMax(key string, count ...int64) *redis.ZSliceCmd {
	return r.client.ZPopMax(key, count...)
}

func (r *Redis) ZPopMin(key string, count ...int64) *redis.ZSliceCmd {
	return r.client.ZPopMin(key, count...)
}

func (r *Redis) ZRange(key string, start, stop int64) *redis.StringSliceCmd {
	return r.client.ZRange(key, start, stop)
}

func (r *Redis) ZRangeWithScores(key string, start, stop int64) *redis.ZSliceCmd {
	return r.client.ZRangeWithScores(key, start, stop)
}

func (r *Redis) ZRangeByScore(key string, opt redis.ZRangeBy) *redis.StringSliceCmd {
	return r.client.ZRangeByScore(key, opt)
}

func (r *Redis) ZRangeByLex(key string, opt redis.ZRangeBy) *redis.StringSliceCmd {
	return r.client.ZRangeByLex(key, opt)
}

func (r *Redis) ZRangeByScoreWithScores(key string, opt redis.ZRangeBy) *redis.ZSliceCmd {
	return r.client.ZRangeByScoreWithScores(key, opt)
}

func (r *Redis) ZRank(key, member string) *redis.IntCmd {
	return r.client.ZRank(key, member)
}

func (r *Redis) ZRem(key string, members ...interface{}) *redis.IntCmd {
	return r.client.ZRem(key, members)
}

func (r *Redis) ZRemRangeByRank(key string, start, stop int64) *redis.IntCmd {
	return r.client.ZRemRangeByRank(key, start, stop)
}

func (r *Redis) ZRemRangeByScore(key, min, max string) *redis.IntCmd {
	return r.client.ZRemRangeByScore(key, min, max)
}

func (r *Redis) ZRemRangeByLex(key, min, max string) *redis.IntCmd {
	return r.client.ZRemRangeByLex(key, min, max)
}

func (r *Redis) ZRevRange(key string, start, stop int64) *redis.StringSliceCmd {
	return r.client.ZRevRange(key, start, stop)
}

func (r *Redis) ZRevRangeWithScores(key string, start, stop int64) *redis.ZSliceCmd {
	return r.client.ZRevRangeWithScores(key, start, stop)
}

func (r *Redis) ZRevRangeByScore(key string, opt redis.ZRangeBy) *redis.StringSliceCmd {
	return r.client.ZRevRangeByScore(key, opt)
}

func (r *Redis) ZRevRangeByLex(key string, opt redis.ZRangeBy) *redis.StringSliceCmd {
	return r.client.ZRevRangeByLex(key, opt)
}

func (r *Redis) ZRevRangeByScoreWithScores(key string, opt redis.ZRangeBy) *redis.ZSliceCmd {
	return r.client.ZRevRangeByScoreWithScores(key, opt)
}

func (r *Redis) ZRevRank(key, member string) *redis.IntCmd {
	return r.client.ZRevRank(key, member)
}

func (r *Redis) ZScore(key, member string) *redis.FloatCmd {
	return r.client.ZScore(key, member)
}

func (r *Redis) ZUnionStore(dest string, store redis.ZStore, keys ...string) *redis.IntCmd {
	return r.client.ZUnionStore(dest, store, keys...)
}

func (r *Redis) PFAdd(key string, els ...interface{}) *redis.IntCmd {
	return r.client.PFAdd(key, els...)
}

func (r *Redis) PFCount(keys ...string) *redis.IntCmd {
	return r.client.PFCount(keys...)
}

func (r *Redis) PFMerge(dest string, keys ...string) *redis.StatusCmd {
	return r.client.PFMerge(dest, keys...)
}

func (r *Redis) BgRewriteAOF() *redis.StatusCmd {
	return r.client.BgRewriteAOF()
}

func (r *Redis) BgSave() *redis.StatusCmd {
	return r.client.BgSave()
}

func (r *Redis) ClientKill(ipPort string) *redis.StatusCmd {
	return r.client.ClientKill(ipPort)
}

func (r *Redis) ClientKillByFilter(keys ...string) *redis.IntCmd {
	return r.client.ClientKillByFilter(keys...)
}

func (r *Redis) ClientList() *redis.StringCmd {
	return r.client.ClientList()
}

func (r *Redis) ClientPause(dur time.Duration) *redis.BoolCmd {
	return r.client.ClientPause(dur)
}

func (r *Redis) ClientID() *redis.IntCmd {
	return r.client.ClientID()
}

func (r *Redis) ConfigGet(parameter string) *redis.SliceCmd {
	return r.client.ConfigGet(parameter)
}

func (r *Redis) ConfigResetStat() *redis.StatusCmd {
	return r.client.ConfigResetStat()
}

func (r *Redis) ConfigSet(parameter, value string) *redis.StatusCmd {
	return r.client.ConfigSet(parameter, value)
}

func (r *Redis) ConfigRewrite() *redis.StatusCmd {
	return r.client.ConfigRewrite()
}

func (r *Redis) DBSize() *redis.IntCmd {
	return r.client.DBSize()
}

func (r *Redis) FlushAll() *redis.StatusCmd {
	return r.client.FlushAll()
}

func (r *Redis) FlushAllAsync() *redis.StatusCmd {
	return r.client.FlushAllAsync()
}

func (r *Redis) FlushDB() *redis.StatusCmd {
	return r.client.FlushDB()
}

func (r *Redis) FlushDBAsync() *redis.StatusCmd {
	return r.client.FlushDBAsync()
}

func (r *Redis) Info(section ...string) *redis.StringCmd {
	return r.client.Info(section...)
}

func (r *Redis) LastSave() *redis.IntCmd {
	return r.client.LastSave()
}

func (r *Redis) Save() *redis.StatusCmd {
	return r.client.Save()
}

func (r *Redis) Shutdown() *redis.StatusCmd {
	return r.client.Shutdown()
}

func (r *Redis) ShutdownSave() *redis.StatusCmd {
	return r.client.ShutdownSave()
}

func (r *Redis) ShutdownNoSave() *redis.StatusCmd {
	return r.client.ShutdownNoSave()
}

func (r *Redis) SlaveOf(host, port string) *redis.StatusCmd {
	return r.client.SlaveOf(host, port)
}

func (r *Redis) Time() *redis.TimeCmd {
	return r.client.Time()
}

func (r *Redis) Eval(script string, keys []string, args ...interface{}) *redis.Cmd {
	return r.client.Eval(script, keys, args...)
}

func (r *Redis) EvalSha(sha1 string, keys []string, args ...interface{}) *redis.Cmd {
	return r.client.EvalSha(sha1, keys, args...)
}

func (r *Redis) ScriptExists(hashes ...string) *redis.BoolSliceCmd {
	return r.client.ScriptExists(hashes...)
}

func (r *Redis) ScriptFlush() *redis.StatusCmd {
	return r.client.ScriptFlush()
}

func (r *Redis) ScriptKill() *redis.StatusCmd {
	return r.client.ScriptKill()
}

func (r *Redis) ScriptLoad(script string) *redis.StringCmd {
	return r.client.ScriptLoad(script)
}

func (r *Redis) DebugObject(key string) *redis.StringCmd {
	return r.client.DebugObject(key)
}

func (r *Redis) Publish(channel string, message interface{}) *redis.IntCmd {
	return r.client.Publish(channel, message)
}

func (r *Redis) PubSubChannels(pattern string) *redis.StringSliceCmd {
	return r.client.PubSubChannels(pattern)
}

func (r *Redis) PubSubNumSub(channels ...string) *redis.StringIntMapCmd {
	return r.client.PubSubNumSub(channels...)
}

func (r *Redis) PubSubNumPat() *redis.IntCmd {
	return r.client.PubSubNumPat()
}

func (r *Redis) ClusterSlots() *redis.ClusterSlotsCmd {
	return r.client.ClusterSlots()
}

func (r *Redis) ClusterNodes() *redis.StringCmd {
	return r.client.ClusterNodes()
}

func (r *Redis) ClusterMeet(host, port string) *redis.StatusCmd {
	return r.client.ClusterMeet(host, port)
}

func (r *Redis) ClusterForget(nodeID string) *redis.StatusCmd {
	return r.client.ClusterForget(nodeID)
}

func (r *Redis) ClusterReplicate(nodeID string) *redis.StatusCmd {
	return r.client.ClusterReplicate(nodeID)
}

func (r *Redis) ClusterResetSoft() *redis.StatusCmd {
	return r.client.ClusterResetSoft()
}

func (r *Redis) ClusterResetHard() *redis.StatusCmd {
	return r.client.ClusterResetHard()
}

func (r *Redis) ClusterInfo() *redis.StringCmd {
	return r.client.ClusterInfo()
}

func (r *Redis) ClusterKeySlot(key string) *redis.IntCmd {
	return r.client.ClusterKeySlot(key)
}

func (r *Redis) ClusterGetKeysInSlot(slot int, count int) *redis.StringSliceCmd {
	return r.client.ClusterGetKeysInSlot(slot, count)
}

func (r *Redis) ClusterCountFailureReports(nodeID string) *redis.IntCmd {
	return r.client.ClusterCountFailureReports(nodeID)
}

func (r *Redis) ClusterCountKeysInSlot(slot int) *redis.IntCmd {
	return r.client.ClusterCountKeysInSlot(slot)
}

func (r *Redis) ClusterDelSlots(slots ...int) *redis.StatusCmd {
	return r.client.ClusterDelSlots(slots...)
}

func (r *Redis) ClusterDelSlotsRange(min, max int) *redis.StatusCmd {
	return r.client.ClusterDelSlotsRange(min, max)
}

func (r *Redis) ClusterSaveConfig() *redis.StatusCmd {
	return r.client.ClusterSaveConfig()
}

func (r *Redis) ClusterSlaves(nodeID string) *redis.StringSliceCmd {
	return r.client.ClusterSlaves(nodeID)
}

func (r *Redis) ClusterFailover() *redis.StatusCmd {
	return r.client.ClusterFailover()
}

func (r *Redis) ClusterAddSlots(slots ...int) *redis.StatusCmd {
	return r.client.ClusterAddSlots(slots...)
}

func (r *Redis) ClusterAddSlotsRange(min, max int) *redis.StatusCmd {
	return r.client.ClusterAddSlotsRange(min, max)
}

func (r *Redis) GeoAdd(key string, geoLocation ...*redis.GeoLocation) *redis.IntCmd {
	return r.client.GeoAdd(key, geoLocation...)
}

func (r *Redis) GeoPos(key string, members ...string) *redis.GeoPosCmd {
	return r.client.GeoPos(key, members...)
}

func (r *Redis) GeoRadius(key string, longitude, latitude float64, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	return r.client.GeoRadius(key, longitude, latitude, query)
}

func (r *Redis) GeoRadiusRO(key string, longitude, latitude float64, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	return r.client.GeoRadiusRO(key, longitude, latitude, query)
}

func (r *Redis) GeoRadiusByMember(key, member string, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	return r.client.GeoRadiusByMember(key, member, query)
}

func (r *Redis) GeoRadiusByMemberRO(key, member string, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	return r.client.GeoRadiusByMemberRO(key, member, query)
}

func (r *Redis) GeoDist(key string, member1, member2, unit string) *redis.FloatCmd {
	return r.client.GeoDist(key, member1, member2, unit)
}

func (r *Redis) GeoHash(key string, members ...string) *redis.StringSliceCmd {
	return r.client.GeoHash(key, members...)
}

func (r *Redis) ReadOnly() *redis.StatusCmd {
	return r.client.ReadOnly()
}

func (r *Redis) ReadWrite() *redis.StatusCmd {
	return r.client.ReadWrite()
}

func (r *Redis) MemoryUsage(key string, samples ...int) *redis.IntCmd {
	return r.client.MemoryUsage(key, samples...)
}

func fromClient(client *redis.Client) RedisIOCInterface {
	return autowire.GetProxyFunction()(&Redis{
		client: client,
	}).(RedisIOCInterface)
}

func (r *Redis) Context() context.Context {
	return r.client.Context()
}

func (r *Redis) WithContext(ctx context.Context) RedisIOCInterface {
	return fromClient(r.client.WithContext(ctx))
}

// Options returns read-only Options that were used to create the client.
func (r *Redis) Options() *redis.Options {
	return r.client.Options()
}

func (r *Redis) SetLimiter(l redis.Limiter) RedisIOCInterface {
	return fromClient(r.client.SetLimiter(l))
}

// PoolStats returns connection pool stats.
func (r *Redis) PoolStats() *redis.PoolStats {
	return r.client.PoolStats()
}

func (r *Redis) Pipelined(fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	return r.client.Pipelined(fn)
}

func (r *Redis) Pipeline() redis.Pipeliner {
	return r.client.Pipeline()
}

func (r *Redis) TxPipelined(fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	return r.client.TxPipelined(fn)
}

// TxPipeline acts like Pipeline, but wraps queued commands with MULTI/EXEC.
func (r *Redis) TxPipeline() redis.Pipeliner {
	return r.client.TxPipeline()
}

// Subscribe subscribes the client to the specified channels.
// Channels can be omitted to create empty subscription.
// Note that this method does not wait on a response from Redis, so the
// subscription may not be active immediately. To force the connection to wait,
// you may call the Receive() method on the returned *PubSub like so:
//
//    sub := client.Subscribe(queryResp)
//    iface, err := sub.Receive()
//    if err != nil {
//        // handle error
//    }
//
//    // Should be *Subscription, but others are possible if other actions have been
//    // taken on sub since it was created.
//    switch iface.(type) {
//    case *Subscription:
//        // subscribe succeeded
//    case *Message:
//        // received first message
//    case *Pong:
//        // pong received
//    default:
//        // handle error
//    }
//
//    ch := sub.Channel()
func (r *Redis) Subscribe(channels ...string) *redis.PubSub {
	return r.client.Subscribe(channels...)
}

// PSubscribe subscribes the client to the given patterns.
// Patterns can be omitted to create empty subscription.
func (r *Redis) PSubscribe(channels ...string) *redis.PubSub {
	return r.client.PSubscribe(channels...)
}
