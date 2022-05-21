package config

import (
	"fmt"
	"testing"
)

func TestMergeMap(t *testing.T) {
	type args struct {
		dst              AnyMap
		src              AnyMap
		depths           []uint8
		wantPanic        bool
		wantPanicContent string
	}
	tests := []struct {
		name string
		args args
		want AnyMap
	}{
		{
			name: "Test MergeMap()-panic",
			args: args{
				dst:              make(AnyMap),
				src:              AnyMap{},
				depths:           []uint8{32},
				wantPanic:        true,
				wantPanicContent: "[Config] expect depth too deep: [32]",
			},
			want: AnyMap{},
		},
		{
			name: "Test MergeMap()",
			args: args{
				dst: make(AnyMap),
				src: AnyMap{
					"hello": "world",
					"Go":    2009,
				},
			},
			want: AnyMap{
				"hello": "world",
				"Go":    2009,
			},
		},
		{
			name: "Test MergeMap()",
			args: args{
				dst: AnyMap{
					"Java": 1995,
				},
				src: AnyMap{
					"hello": "world",
					"Go":    2009,
				},
			},
			want: AnyMap{
				"Java":  1995,
				"hello": "world",
				"Go":    2009,
			},
		},
		{
			name: "Test MergeMap()",
			args: args{
				dst: AnyMap{
					"Java":  1995,
					"hello": "world",
				},
				src: AnyMap{
					"hello": "world",
					"Go":    2009,
				},
			},
			want: AnyMap{
				"Java":  1995,
				"hello": "world",
				"Go":    2009,
			},
		},
		{
			name: "Test MergeMap()-depth_8",
			args: args{
				dst: AnyMap{
					"Java":  1995,
					"Go":    2009,
					"hello": "world",
					"key1": AnyMap{
						"Java":  1995,
						"hello": "world",
						"key2": AnyMap{
							"Java":  1995,
							"hello": "world",
							"key3": AnyMap{
								"Java":  1995,
								"hello": "world",
								"key4": AnyMap{
									"Java":  1995,
									"hello": "world",
									"key5": AnyMap{
										"Java":  1995,
										"hello": "world",
										"key6": AnyMap{
											"Java":  1995,
											"hello": "world",
											"key7": AnyMap{
												"Java":  1995,
												"hello": "world",
												"key8": AnyMap{
													"Java":  1995,
													"hello": "world",
												},
											},
										},
									},
								},
							},
						},
					},
				},
				src: AnyMap{
					"Java":  1995,
					"hello": "world",
					"key1": AnyMap{
						"Java":  1995,
						"hello": "world",
						"key2": AnyMap{
							"Java":  1995,
							"hello": "world",
							"key3": AnyMap{
								"Java":  1995,
								"hello": "world",
								"key4": AnyMap{
									"Java":  1995,
									"hello": "world",
									"key5": AnyMap{
										"Java":  1995,
										"hello": "world",
										"key6": AnyMap{
											"Java":  1995,
											"hello": "world",
											"key7": AnyMap{
												"Java":  1995,
												"hello": "world",
												"key8": AnyMap{
													"Java":  1995,
													"hello": "world",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: AnyMap{
				"Java":  1995,
				"Go":    2009,
				"hello": "world",
				"key1": AnyMap{
					"Java":  1995,
					"hello": "world",
					"key2": AnyMap{
						"Java":  1995,
						"hello": "world",
						"key3": AnyMap{
							"Java":  1995,
							"hello": "world",
							"key4": AnyMap{
								"Java":  1995,
								"hello": "world",
								"key5": AnyMap{
									"Java":  1995,
									"hello": "world",
									"key6": AnyMap{
										"Java":  1995,
										"hello": "world",
										"key7": AnyMap{
											"Java":  1995,
											"hello": "world",
											"key8": AnyMap{
												"Java":  1995,
												"hello": "world",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Test MergeMap()-depth16",
			args: args{
				dst:       depth16,
				src:       depth16,
				depths:    []uint8{16},
				wantPanic: true,
			},
			want: depth16,
		},
		{
			name: "Test MergeMap()-depth17",
			args: args{
				dst:              depth17,
				src:              depth17,
				depths:           []uint8{16},
				wantPanic:        true,
				wantPanicContent: "[Config] recursion too deep: [17]",
			},
			want: AnyMap{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.wantPanic {
				defer func() {
					if err := recover(); err != nil {
						if fmt.Sprintf("%v", err) != tt.args.wantPanicContent {
							panic(err)
						}
					}
				}()
			}

			am := MergeMap(tt.args.dst, tt.args.src, tt.args.depths...)
			t.Logf("am is:%v", am)
		})
	}
}

var depth16 = AnyMap{
	"Java":  1995,
	"Go":    2009,
	"hello": "world",
	"key1": AnyMap{
		"Java":  1995,
		"hello": "world",
		"key2": AnyMap{
			"Java":  1995,
			"hello": "world",
			"key3": AnyMap{
				"Java":  1995,
				"hello": "world",
				"key4": AnyMap{
					"Java":  1995,
					"hello": "world",
					"key5": AnyMap{
						"Java":  1995,
						"hello": "world",
						"key6": AnyMap{
							"Java":  1995,
							"hello": "world",
							"key7": AnyMap{
								"Java":  1995,
								"hello": "world",
								"key8": AnyMap{
									"Java":  1995,
									"hello": "world",
									"key9": AnyMap{
										"Java":  1995,
										"hello": "world",
										"key10": AnyMap{
											"Java":  1995,
											"hello": "world",
											"key11": AnyMap{
												"Java":  1995,
												"hello": "world",
												"key12": AnyMap{
													"Java":  1995,
													"hello": "world",
													"key13": AnyMap{
														"Java":  1995,
														"hello": "world",
														"key14": AnyMap{
															"Java":  1995,
															"hello": "world",
															"key15": AnyMap{
																"Java":  1995,
																"hello": "world",
																"key16": AnyMap{
																	"Java":  1995,
																	"hello": "world",
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
}

var depth17 = AnyMap{
	"Java":  1995,
	"Go":    2009,
	"hello": "world",
	"key1": AnyMap{
		"Java":  1995,
		"hello": "world",
		"key2": AnyMap{
			"Java":  1995,
			"hello": "world",
			"key3": AnyMap{
				"Java":  1995,
				"hello": "world",
				"key4": AnyMap{
					"Java":  1995,
					"hello": "world",
					"key5": AnyMap{
						"Java":  1995,
						"hello": "world",
						"key6": AnyMap{
							"Java":  1995,
							"hello": "world",
							"key7": AnyMap{
								"Java":  1995,
								"hello": "world",
								"key8": AnyMap{
									"Java":  1995,
									"hello": "world",
									"key9": AnyMap{
										"Java":  1995,
										"hello": "world",
										"key10": AnyMap{
											"Java":  1995,
											"hello": "world",
											"key11": AnyMap{
												"Java":  1995,
												"hello": "world",
												"key12": AnyMap{
													"Java":  1995,
													"hello": "world",
													"key13": AnyMap{
														"Java":  1995,
														"hello": "world",
														"key14": AnyMap{
															"Java":  1995,
															"hello": "world",
															"key15": AnyMap{
																"Java":  1995,
																"hello": "world",
																"key16": AnyMap{
																	"Java":  1995,
																	"hello": "world",
																	"key17": AnyMap{
																		"Java":  1995,
																		"hello": "world",
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
}
