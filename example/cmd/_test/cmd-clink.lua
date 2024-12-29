local function example_completion(word, word_index, line_state, match_builder)
  local compline = string.sub(line_state:getline(), 1, line_state:getcursor())

  local output = io.popen("env CARAPACE_COMPLINE=" .. string.format("%q", compline) .. " example _carapace cmd-clink ''"):read("*a")
  for line in string.gmatch(output, '[^\r\n]+') do
    match_builder:addmatch(string.gsub(line, '\t.*', ""))
  end
  return true
end

clink.argmatcher("example"):addarg({example_completion}):loop(1)

