package reminder

const (
	getReminderByID        = `SELECT id, text, user_id, remind_at, remind_original, snoozed FROM  reminders WHERE id = $1`
	insertReminder         = `INSERT INTO reminders(text,user_id,remind_at,remind_original) VALUES($1,$2,$3,$4)`
	getRemindersTodayByUID = `
								SELECT id, text, user_id, remind_at, remind_original, snoozed 
								FROM reminders 
								WHERE user_id = $1 AND remind_at::date < current_date + 1
							`
)
