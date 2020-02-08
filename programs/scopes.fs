let a = "global"
let b = "global"
let c = "global"
{
  print a
  b = "first"
  print b
  print c

  {
    print a
    print b
    c = "second"
    print c

    let d = "inner let"
  }

  print a
  print b
  print c
}

print a; print b; print c
