

past = ""
total = 0

while True:
    packet = input("> ")
    if not packet:
        break

    pastSet = set(past)
    if set(packet).issubset(pastSet):
        print("good")
    else:
        print(pow((len(pastSet)/(len(pastSet)+1)), total))
    past += packet
    total += len(packet)
