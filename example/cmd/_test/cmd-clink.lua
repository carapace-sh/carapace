local function example_completion(word, word_index, line_state, match_builder)
  local compline = string.sub(line_state:getline(), 1, line_state:getcursor())

  local output = io.popen("env CARAPACE_COMPLINE=" .. string.format("%q", compline) .. " example _carapace cmd-clink \"\""):read("*a")
  for line in string.gmatch(output, '[^\r\n]+') do
    local matches = {}
    for m in string.gmatch(line, '[^\t]+') do
      table.insert(matches, m)
    end
    match_builder:addmatch({
      match = matches[1],
      display = matches[2],
      description = matches[3],
      type = "word",
      appendchar = matches[4],
      suppressappend = false
    })
  end
  return true
end

clink.argmatcher("example"):addarg({nowordbreakchars="'&backprime;+;,", example_completion}):loop(1)

