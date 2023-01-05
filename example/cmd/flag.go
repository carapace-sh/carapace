package cmd

import (
	"net"
	"time"

	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

var flagCmd = &cobra.Command{
	Use:     "flag",
	Short:   "flag example",
	GroupID: "main",
	Run:     func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(flagCmd)

	flagCmd.Flags().Bool("Bool", false, "Bool")
	flagCmd.Flags().BoolSlice("BoolSlice", []bool{}, "BoolSlice")
	flagCmd.Flags().BytesBase64("BytesBase64", []byte{}, "BytesBase64")
	flagCmd.Flags().BytesHex("BytesHex", []byte{}, "BytesHex")
	flagCmd.Flags().Count("Count", "Count")
	flagCmd.Flags().Duration("Duration", 0, "Duration")
	flagCmd.Flags().DurationSlice("DurationSlice", []time.Duration{}, "DurationSlice")
	flagCmd.Flags().Float32("Float32P", 0, "Float32P")
	flagCmd.Flags().Float32Slice("Float32Slice", []float32{}, "Float32Slice")
	flagCmd.Flags().Float64("Float64P", 0, "Float64P")
	flagCmd.Flags().Float64Slice("Float64Slice", []float64{}, "Float64Slice")
	flagCmd.Flags().Int16("Int16", 0, "Int16")
	flagCmd.Flags().Int32("Int32", 0, "Int32")
	flagCmd.Flags().Int32Slice("Int32Slice", []int32{}, "Int32Slice")
	flagCmd.Flags().Int64("Int64", 0, "Int64")
	flagCmd.Flags().Int64Slice("Int64Slice", []int64{}, "Int64Slice")
	flagCmd.Flags().Int8("Int8", 0, "Int8")
	flagCmd.Flags().Int("Int", 0, "Int")
	flagCmd.Flags().IntSlice("IntSlice", []int{}, "IntSlice")
	flagCmd.Flags().IPMask("IPMask", net.IPMask{}, "IPMask")
	flagCmd.Flags().IP("IP", net.IP{}, "IP")
	flagCmd.Flags().IPNet("IPNet", net.IPNet{}, "IPNet")
	flagCmd.Flags().IPSlice("IPSlice", []net.IP{}, "IPSlice")
	flagCmd.Flags().StringArray("StringArray", []string{}, "StringArray")
	flagCmd.Flags().String("String", "", "String")
	flagCmd.Flags().StringSlice("StringSlice", []string{}, "StringSlice")
	flagCmd.Flags().StringToInt64("StringToInt64", map[string]int64{}, "StringToInt64")
	flagCmd.Flags().StringToInt("StringToInt", map[string]int{}, "StringToInt")
	flagCmd.Flags().StringToString("StringToString", map[string]string{}, "StringToString")
	flagCmd.Flags().Uint16("Uint16", 0, "Uint16")
	flagCmd.Flags().Uint32("Uint32", 0, "Uint32")
	flagCmd.Flags().Uint64("Uint64", 0, "Uint64")
	flagCmd.Flags().Uint8("Uint8", 0, "Uint8")
	flagCmd.Flags().Uint("Uint", 0, "Uint")
	flagCmd.Flags().UintSlice("UintSlice", []uint{}, "UintSlice")

	flagCmd.Flags().Bool("optarg", false, "test optarg variant (must be second arg on command line to work)") // TODO quick&dirty toggle for now
	carapace.Gen(rootCmd).PreRun(func(cmd *cobra.Command, args []string) {
		if len(args) < 2 || args[1] != "--optarg" {
			return
		}

		// TODO set correct default values
		flagCmd.Flag("Bool").NoOptDefVal = " "
		flagCmd.Flag("BoolSlice").NoOptDefVal = " "
		flagCmd.Flag("BytesBase64").NoOptDefVal = " "
		flagCmd.Flag("BytesHex").NoOptDefVal = " "
		flagCmd.Flag("Count").NoOptDefVal = " "
		flagCmd.Flag("Duration").NoOptDefVal = " "
		flagCmd.Flag("DurationSlice").NoOptDefVal = " "
		flagCmd.Flag("Float32P").NoOptDefVal = " "
		flagCmd.Flag("Float32Slice").NoOptDefVal = " "
		flagCmd.Flag("Float64P").NoOptDefVal = " "
		flagCmd.Flag("Float64Slice").NoOptDefVal = " "
		flagCmd.Flag("Int16").NoOptDefVal = " "
		flagCmd.Flag("Int32").NoOptDefVal = " "
		flagCmd.Flag("Int32Slice").NoOptDefVal = " "
		flagCmd.Flag("Int64").NoOptDefVal = " "
		flagCmd.Flag("Int64Slice").NoOptDefVal = " "
		flagCmd.Flag("Int8").NoOptDefVal = " "
		flagCmd.Flag("Int").NoOptDefVal = " "
		flagCmd.Flag("IntSlice").NoOptDefVal = " "
		flagCmd.Flag("IPMask").NoOptDefVal = " "
		flagCmd.Flag("IP").NoOptDefVal = " "
		flagCmd.Flag("IPNet").NoOptDefVal = " "
		flagCmd.Flag("IPSlice").NoOptDefVal = " "
		flagCmd.Flag("StringArray").NoOptDefVal = " "
		flagCmd.Flag("String").NoOptDefVal = " "
		flagCmd.Flag("StringSlice").NoOptDefVal = " "
		flagCmd.Flag("StringToInt64").NoOptDefVal = " "
		flagCmd.Flag("StringToInt").NoOptDefVal = " "
		flagCmd.Flag("StringToString").NoOptDefVal = " "
		flagCmd.Flag("Uint16").NoOptDefVal = " "
		flagCmd.Flag("Uint32").NoOptDefVal = " "
		flagCmd.Flag("Uint64").NoOptDefVal = " "
		flagCmd.Flag("Uint8").NoOptDefVal = " "
		flagCmd.Flag("Uint").NoOptDefVal = " "
		flagCmd.Flag("UintSlice").NoOptDefVal = " "
	})

	carapace.Gen(flagCmd).FlagCompletion(carapace.ActionMap{
		"Bool":           carapace.ActionValues("true", "false"),
		"BoolSlice":      carapace.ActionValues("true", "false"),
		"BytesBase64":    carapace.ActionValues("MQo=", "Mgo=", "Mwo="),
		"BytesHex":       carapace.ActionValues("01", "02", "03"),
		"Count":          carapace.ActionValues(),
		"Duration":       carapace.ActionValues("1h", "2m", "3s"),
		"DurationSlice":  carapace.ActionValues("1h", "2m", "3s"),
		"Float32P":       carapace.ActionValues("1", "2", "3"),
		"Float32Slice":   carapace.ActionValues("1", "2", "3"),
		"Float64P":       carapace.ActionValues("1", "2", "3"),
		"Float64Slice":   carapace.ActionValues("1", "2", "3"),
		"Int16":          carapace.ActionValues("1", "2", "3"),
		"Int32":          carapace.ActionValues("1", "2", "3"),
		"Int32Slice":     carapace.ActionValues("1", "2", "3"),
		"Int64":          carapace.ActionValues("1", "2", "3"),
		"Int64Slice":     carapace.ActionValues("1", "2", "3"),
		"Int8":           carapace.ActionValues("1", "2", "3"),
		"Int":            carapace.ActionValues("1", "2", "3"),
		"IntSlice":       carapace.ActionValues("1", "2", "3"),
		"IPMask":         carapace.ActionValues("0.0.0.1", "0.0.0.2", "0.0.0.3"),
		"IP":             carapace.ActionValues("0.0.0.1", "0.0.0.2", "0.0.0.3"),
		"IPNet":          carapace.ActionValues("0.0.0.1/0", "0.0.0.2/0", "0.0.0.3/0"),
		"IPSlice":        carapace.ActionValues("0.0.0.1", "0.0.0.2", "0.0.0.3"),
		"StringArray":    carapace.ActionValues("1", "2", "3"),
		"String":         carapace.ActionValues("1", "2", "3"),
		"StringSlice":    carapace.ActionValues("1", "2", "3"),
		"StringToInt64":  carapace.ActionValues("a=1", "b=2", "c=3"),
		"StringToInt":    carapace.ActionValues("a=1", "b=2", "c=3"),
		"StringToString": carapace.ActionValues("a=1", "b=2", "c=3"),
		"Uint16":         carapace.ActionValues("1", "2", "3"),
		"Uint32":         carapace.ActionValues("1", "2", "3"),
		"Uint64":         carapace.ActionValues("1", "2", "3"),
		"Uint8":          carapace.ActionValues("1", "2", "3"),
		"Uint":           carapace.ActionValues("1", "2", "3"),
		"UintSlice":      carapace.ActionValues("1", "2", "3"),
	})
}
