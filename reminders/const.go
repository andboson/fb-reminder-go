package reminders

const (
	setSnoozeReminderByID  = `UPDATE reminders SET snoozed = true WHERE id = $1`
	deleteReminderByID     = `DELETE FROM reminders WHERE id = $1`
	deleteReminderByUserID = `DELETE FROM reminders WHERE user_id = $1`
	getReminderByID        = `SELECT id, text, user_id, remind_at, remind_original, snoozed FROM  reminders WHERE id = $1`
	insertReminder         = `INSERT INTO reminders(text,user_id,remind_at,remind_original) VALUES($1,$2,$3,$4)`
	getRemindersTodayByUID = `
								SELECT id, text, user_id, remind_at, remind_original, snoozed 
								FROM reminders 
								WHERE user_id = $1 AND remind_at::date < current_date + 1
							`
	getExpired = `
					SELECT id, text, user_id, remind_at, remind_original, snoozed 
					FROM reminders 
					WHERE 
						(remind_at < now() AND snoozed = false) 
					OR 
						(remind_at + interval '%int%' < now() AND snoozed = true)
				`
)
