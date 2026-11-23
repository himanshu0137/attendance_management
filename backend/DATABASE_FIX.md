# Database Connection Fix Guide

## Problem Analysis

Based on the error message:
```
2026/11/23 05:36:36 ERROR Failed to connect to database error="failed to ping database: failed to connect to `user=carbon database=`: /tmp/.s.PGSQL.5432 (/tmp): dial error: dial unix /tmp/.s.PGSQL.5432: connect: no such file or directory"
```

The issue is that the PostgreSQL connection string is not being properly parsed. The system is trying to use a Unix socket connection instead of TCP connection.

## Root Cause

The environment variable `ATTENDANCE_DB` is either:
1. Not set properly 
2. Set to an empty value
3. Set to a malformed connection string
4. Not accessible from the Go application context

## Solution Steps

### 1. Verify Database Connection String Format

The correct PostgreSQL connection string should look like:
```
postgresql://username:password@host:port/database_name
```

For Docker containers, it should be:
```
postgresql://attendance:attendance@localhost:5432/attendance
```

### 2. Check Environment Variable in Go Code

Let's create a debugging version of the main.go to see what's actually being loaded:

<tool_call>
<function=write>
<parameter=path>
/home/carbon/prog/attendance_management/backend/debug_main.go