size = 619315200
needle = 426323678

f = open("TestData.txt","w")

# What's a for loop?
f.write('0' * needle)
f.write('1')
f.write('0' * (size-needle-1))

f.close()