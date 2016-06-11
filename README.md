# mailnotify

Provide simple notifications via pushover, usable in your .procmailrc:

	COMSAT=no
	SUBJECT=`formail -xSubject:`    # regular field
	MSGID=`formail -xMessage-ID:`    # regular field
	FROM=`formail -rt -xTo:`        # special case
	
	:0ic
	| mail_notify

and mail_notify a script to keep you pushover credentials:

	#!/bin/sh
	$HOME/go/bin/mailnotify api_key user_key
