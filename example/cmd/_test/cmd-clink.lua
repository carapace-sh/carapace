local function example_completion(word, word_index, line_state, match_builder)
  args = { "example", "_carapace", "fish", "\"\"" }
  for i = 2,word_index-1,1 do
    table.insert(args, string.format("%q" ,line_state:getword(i)))
  end

  -- table.insert(args, string.format("%q", word))  
  local exploded = string.explode(line_state:getline() .. "a") 
  word = string.gsub(exploded[#exploded], 'a$', "") 
  table.insert(args, string.format("%q", word)) 

  output = io.popen(table.concat(args, " ")):read("*a")
  for line in string.gmatch(output, '[^\r\n]+') do
    -- match_builder:addmatch(line)
    match_builder:addmatch(string.gsub(line, '\t.*', ""))
  end

  return true
end

clink.argmatcher("example"):addarg({example_completion}):loop(1)

