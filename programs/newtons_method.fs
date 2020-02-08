// Find the square root of x
let x = 8
let y = x

while x*x - y > 0.0000000000000001 then
    x = (x + y/x)/2


print "The square root of " + y + " is " + x + "."