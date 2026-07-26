import 'package:flutter_test/flutter_test.dart';
import 'package:app/main.dart';

void main() {
  testWidgets('App loads without crashing', (WidgetTester tester) async {
    await tester.pumpWidget(const MiTramiteApp());
    expect(find.byType(MiTramiteApp), findsOneWidget);
  });
}
