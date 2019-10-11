import yahoo_fin.options as yf, pandas as pd, datetime as dt, calendar, sys, time

if __name__ == '__main__':
    c, y, m = calendar.Calendar(firstweekday=calendar.SUNDAY), int(dt.datetime.strftime(dt.datetime.now(), '%Y')), int(dt.datetime.strftime(dt.datetime.now(), '%m'))
    if dt.datetime.now().day > 21: m += 1
    month_cals = [c.monthdatescalendar(y, n) for n in range(m, m + 4)]
    opex = [dt.datetime.strftime(day, '%Y-%m-%d') for mo in month_cals for week in mo for day in week if day.weekday() == 4 and 15 <= day.day <= 21]

    for n in range(len(opex)):
        try:
            time.sleep(1.5)
            
            print(f'\n\n{sys.argv[1].upper()} {opex[n]}\n')

            opex_ps = yf.get_puts(sys.argv[1], opex[n])
            puts = opex_ps[opex_ps['Open Interest'] > 0].sort_values('Volume', ascending=False).reset_index()
            print(puts[['Contract Name', 'Strike', 'Last Price', 'Bid', 'Ask', 'Volume', 'Open Interest', 'Change', '% Change']][0:10])

            del opex_ps, puts
        
        except ValueError as _:
            continue

    # for n in range(len(opex)):
        # print(f'\n\n{sys.argv[1].upper()} {opex[n]}\n')        
        
        # opex_cs = yf.get_calls(sys.argv[1], opex[0])
        # calls = opex_cs[opex_cs['Open Interest'] > 0].sort_values('Volume', ascending=False).reset_index()
        # print(calls[['Contract Name', 'Strike', 'Last Price', 'Bid', 'Ask', 'Volume', 'Open Interest', 'Change', '% Change']][0:10])