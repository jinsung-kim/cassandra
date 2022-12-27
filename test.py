backends = []

nodes=3
rf=2

start=8090

for i in range(nodes*rf):
    backends.append(start)
    start += 1

print(backends)

for i in range(0, nodes*rf, rf):
    for j in range(int(rf % (i + 1))):
        print(j)
    print()