package main

import "testing"

func Test_grep(t *testing.T) {
	type args struct {
		flags       *flags
		nonFlagArgs []string
	}
	type test struct {
		name    string
		args    args
		want    string
		wantErr bool
	}

	// Test case 1 ----------------------------------------------------------------------------
	flags1 := &flags{}
	nonFlagArgs1 := []string{"Aenean", "test_text.txt"}
	args1 := args{flags1, nonFlagArgs1}
	want1 := `test_text.txt:
Aenean commodo ligula eget dolor.
Aenean massa.
Aenean vulputate eleifend tellus.
Aenean leo ligula, porttitor eu, consequat vitae, eleifend ac, enim.
---
`
	test1 := test{"Simple grep by word", args1, want1, false}

	// Test case 2 ----------------------------------------------------------------------------
	flags2 := &flags{}
	nonFlagArgs2 := []string{`[-]?\d[\d,]*[\.]?[\d{2}]*`, "test_text.txt"}
	args2 := args{flags2, nonFlagArgs2}
	want2 := `test_text.txt:
Donec quam felis, ultricies 538 nec, pellentesque eu, pretium quis, sem.
Nulla consequat mas987sa quis enim.
Cras dapibus 1876.
---
`
	test2 := test{"Simple grep by regexp", args2, want2, false}

	// Test case 3 ---------------------------------------------------------------------------------------------
	flags3 := &flags{}
	nonFlagArgs3 := []string{`Aenean`, "test_text.txt", "test_text2.txt"}
	args3 := args{flags3, nonFlagArgs3}
	want3 := `test_text.txt:
Aenean commodo ligula eget dolor.
Aenean massa.
Aenean vulputate eleifend tellus.
Aenean leo ligula, porttitor eu, consequat vitae, eleifend ac, enim.
---
test_text2.txt:
Aenean commodo ligula eget dolor. Aenean massa.
Vivamus elementum semper nisi. Aenean vulputate eleifend tellus.
Aenean leo ligula, porttitor eu, consequat vitae, eleifend ac, enim.
Phasellus viverra nulla ut metus varius laoreet. Quisque rutrum. Aenean imperdiet.
---
`
	test3 := test{"Grep by multiple files", args3, want3, false}

	// Test case 4 ----------------------------------------------------------------------------
	flags4 := &flags{A: 2}
	nonFlagArgs4 := []string{"Aenean", "test_text.txt"}
	args4 := args{flags4, nonFlagArgs4}
	want4 := `test_text.txt:
Aenean commodo ligula eget dolor.
Aenean massa.
Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus.
--
Aenean massa.
Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus.
Donec quam felis, ultricies 538 nec, pellentesque eu, pretium quis, sem.
--
Aenean vulputate eleifend tellus.
Aenean leo ligula, porttitor eu, consequat vitae, eleifend ac, enim.
Aliquam lorem ante, dapibus in, viverra quis, feugiat.
--
Aenean leo ligula, porttitor eu, consequat vitae, eleifend ac, enim.
Aliquam lorem ante, dapibus in, viverra quis, feugiat.
--
---
`
	test4 := test{"Simple grep by word and -A 2", args4, want4, false}

	// Test case 5 ----------------------------------------------------------------------------
	flags5 := &flags{A: 2, B: 1}
	nonFlagArgs5 := []string{"Aenean", "test_text.txt"}
	args5 := args{flags5, nonFlagArgs5}
	want5 := `test_text.txt:
Lorem ipsum dolor sit amet, consectetuer adipiscing elit.
Aenean commodo ligula eget dolor.
Aenean massa.
Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus.
--
Aenean commodo ligula eget dolor.
Aenean massa.
Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus.
Donec quam felis, ultricies 538 nec, pellentesque eu, pretium quis, sem.
--
Vivamus elementum semper nisi.
Aenean vulputate eleifend tellus.
Aenean leo ligula, porttitor eu, consequat vitae, eleifend ac, enim.
Aliquam lorem ante, dapibus in, viverra quis, feugiat.
--
Aenean vulputate eleifend tellus.
Aenean leo ligula, porttitor eu, consequat vitae, eleifend ac, enim.
Aliquam lorem ante, dapibus in, viverra quis, feugiat.
--
---
`
	test5 := test{"Simple grep by word and -A 2 and -B 2", args5, want5, false}

	// Test case 6 ----------------------------------------------------------------------------
	flags6 := &flags{n: true}
	nonFlagArgs6 := []string{"Aenean", "test_text.txt"}
	args6 := args{flags6, nonFlagArgs6}
	want6 := `test_text.txt:
1:Aenean commodo ligula eget dolor.
2:Aenean massa.
12:Aenean vulputate eleifend tellus.
13:Aenean leo ligula, porttitor eu, consequat vitae, eleifend ac, enim.
---
`
	test6 := test{"Simple grep by word with num line", args6, want6, false}

	// Test case 7 ----------------------------------------------------------------------------
	flags7 := &flags{v: true}
	nonFlagArgs7 := []string{"Aenean", "test_text.txt"}
	args7 := args{flags7, nonFlagArgs7}
	want7 := `test_text.txt:
Lorem ipsum dolor sit amet, consectetuer adipiscing elit.
Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus.
Donec quam felis, ultricies 538 nec, pellentesque eu, pretium quis, sem.
Nulla consequat mas987sa quis enim.
Donec pede justo, fringilla vel, aliquet nec, vulputate eget, arcu.
In enim justo, rhoncus ut, imperdiet a, venenatis vitae, justo.
Nullam dictum felis eu pede mollis pretium.
Integer tincidunt.
Cras dapibus 1876.
Vivamus elementum semper nisi.
Aliquam lorem ante, dapibus in, viverra quis, feugiat.
---
`
	test7 := test{"Simple grep by word inverted", args7, want7, false}

	// Test case 8 ----------------------------------------------------------------------------
	flags8 := &flags{c: true}
	nonFlagArgs8 := []string{"Aenean", "test_text.txt"}
	args8 := args{flags8, nonFlagArgs8}
	want8 := `test_text.txt:4
`
	test8 := test{"Simple grep count by word", args8, want8, false}

	// Test case 9 ----------------------------------------------------------------------------
	flags9 := &flags{c: true, v: true}
	nonFlagArgs9 := []string{"Aenean", "test_text.txt"}
	args9 := args{flags9, nonFlagArgs9}
	want9 := `test_text.txt:11
`
	test9 := test{"Simple grep count by word inverted", args9, want9, false}

	// Test case 10 ----------------------------------------------------------------------------
	flags10 := &flags{i: true}
	nonFlagArgs10 := []string{"aEneaN", "test_text.txt"}
	args10 := args{flags10, nonFlagArgs10}
	want10 := `test_text.txt:
Aenean commodo ligula eget dolor.
Aenean massa.
Aenean vulputate eleifend tellus.
Aenean leo ligula, porttitor eu, consequat vitae, eleifend ac, enim.
---
`
	test10 := test{"Simple grep by word ignore case", args10, want10, false}
	// -------------------------------------------------------------------------------------------------

	tests := []test{
		test1,
		test2,
		test3,
		test4,
		test5,
		test6,
		test7,
		test8,
		test9,
		test10,
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := grep(tt.args.flags, tt.args.nonFlagArgs)
			if (err != nil) != tt.wantErr {
				t.Errorf("grep() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("grep() got = %v, want %v", got, tt.want)
			}
		})
	}
}
