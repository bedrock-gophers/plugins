package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func inspectUnsafeWritePacket(path string) error {
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		return err
	}
	for _, declaration := range file.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if !ok || function.Recv != nil || function.Name.Name != "WritePacket" {
			continue
		}
		if function.Type.TypeParams == nil || len(function.Type.TypeParams.List) != 1 {
			return fmt.Errorf("unsafe.WritePacket type parameters changed")
		}
		typeParameter := function.Type.TypeParams.List[0]
		constraint := formatGoExpression(typeParameter.Type)
		if len(typeParameter.Names) != 1 || typeParameter.Names[0].Name != "T" ||
			constraint != "*player.Player | *session.Session" {
			return fmt.Errorf("unsafe.WritePacket target constraint changed: %s", constraint)
		}
		if signature := goFunctionSignature(function); signature != (goSignature{Parameters: "T, packet.Packet"}) {
			return fmt.Errorf("unsafe.WritePacket signature changed: %+v", signature)
		}
		return nil
	}
	return fmt.Errorf("unsafe.WritePacket not found")
}

func generatePlayerWritePacket() []byte {
	var output bytes.Buffer
	output.WriteString(`// Code generated from bedrock-gophers/unsafe WritePacket Go AST. DO NOT EDIT.
#nullable enable
using System;

namespace Dragonfly;

public sealed partial class Player
{
    public void WritePacket(Packet.Packet pk)
    {
        ArgumentNullException.ThrowIfNull(pk);
        PluginBridge.Host.WritePlayerPacket(_invocation, Id, pk.HostHandle());
    }
}
`)
	return output.Bytes()
}
