# SITT
Abbreviation for 'schedule in the terminal', very poetic.

## What's SITT
I like to start my day with some planning. Nothing extreme, just approximate strcuture of the day + some events that are actually scheduled for specific time (if present). While there's probably a bazillion of solutions for such task, I wanted to build my own. And I did. SITT is a cli app for scheduling with syntax close to natural language.

## Usage

Create an entry with following syntax:
```
{NAME} from {TIME} to {TIME}
work from 11:30 to 15:15

{NAME} from {TIME} for {DURATION}
work from 11:30 for 3h
```

Instead of specifying start time for the entry you can use `then' keyword. End time of last entry (last in the schedule, not last one inserted) will be used as start time for you entry.

```
then {NAME} until {TIME}
then work until 13

then {NAME} for {DURATION}
then work for 1h45m
```

### Time/Duration
Specific time should be specified in the following format `{HOUR}:{MINUTE}`, where minutes part is optional. Trailing zeros are also optional. Examples of valid time:

```
07 == 07:00
7 == 07:00
13 == 13:00
13:00 == 13:00
13:5 == 13:05
13:05 == 13:05
```

`now` is a keyword then can be used as time value. It's value corresponds to current system time.

```
work from now for 5h
```

`24 / 24:00` is a valid time.

```
work from now to 24
```

Duration should be specifed in following format `{HOURS}h{MINUTES}m`. Both hours and minutes are optional, but at least one of them should be present. Examples of valid duration:

```
3h
45m
3h45m
```

### Entry Name
Entry name must not contain whitespaces. You can use `_` or other delimiters for multiword names.

`clear` is a special name. Under the hood it inserts entry with an empty name, so syntax described above applies to it.

```
clear from 14:30 to 16:30
```

### Until/To
`until` and `to` are equal in meaning. All of this are valid:
```
then work until 14:00
then work to 14:00
work from 14:00 to 15:00
work from 14:00 until 15:00
```

## Schedule for tomorrow?
Currently only scheduling for today is supported.

## Want to do something fancy with your schedules?
All the data is stored in `$HOMEDIR/sitt-storage` in JSON format. Modifying it is strongly discouraged.
