import datetime
current_date = datetime.datetime.now()
end_date = datetime.timedelta(days=30)
print(f"Current timestamp is {current_date}, end timestamp is {end_date}")
print(current_date + end_date)
