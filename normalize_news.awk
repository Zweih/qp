BEGIN {
  seen_whats_changed = 0
  first_release = 1
}

/^## v[0-9]+\./ {
  if (!first_release) {
    print ""
    print ""
  }

  print
  first_release = 0

  next
}

/^### What's Changed/ {
  seen_whats_changed = 1
  print

  next
}

/^\*\*Full Changelog\*\*:/ {
  if (seen_whats_changed) {
    print ""
  }

  print
  seen_whats_changed = 0

  next
}

{
  print
}
